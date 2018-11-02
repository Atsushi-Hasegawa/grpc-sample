package main
import(
    "log"
    "fmt"
    "net"
    "errors"
    "context"
    "reflect"
	"io/ioutil"
    pb "../../proto"
	"google.golang.org/grpc"
	"github.com/gocraft/dbr"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

type server struct {
}

type Config_Data struct {
	Host     string
	Port     int
	UserName string
	Password string
	Database string
}

type User struct {
	UserId         int64  `db:"user_id"`
	Name           string `db:"name"`
	CreatedDate    string `db:created_date`
	UpdatedDate    string `db:updated_date`
	LastAccessDate string `db:last_access_date`
}

func load_config(filename string) Config_Data {
	buff, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var config Config_Data
	err = yaml.Unmarshal(buff, &config)
	return config
}

func newConnection(file_path string) *dbr.Connection {
	var config = load_config(file_path)
	db, err := dbr.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.UserName, config.Password, config.Host, config.Port, config.Database), nil)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func select_user_info(db *dbr.Connection, id int64) (*User, error) {
	sees := db.NewSession(nil)
    var user *User
    _,err := sees.Select("user.*").
                    From("user").
                    LeftJoin("user_item", "user.user_id = user_item.user_id").
                    Where("user.user_id = ?", id).
                    Load(&user)
    return user, err
}

func (s *server) GetPerson(ctx context.Context, message *pb.GetMessage) (*pb.GetPersonResponse, error) {
	var db = newConnection("../../config/config.yml")
    user,err := select_user_info(db, message.TargetType)
    if err != nil {
        log.Fatal(err)
    }
    response := reflect.ValueOf(user)
    if response.IsNil() {
        return nil, errors.New("")
    }

    return &pb.GetPersonResponse {
        Id: user.UserId,
        Name: user.Name,
        CreatedDate: user.CreatedDate,
        UpdatedDate: user.UpdatedDate,
        LastAccessDate: user.LastAccessDate,
    }, nil
}

func main() {
	listenPort,err := net.Listen("tcp",":19003")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterPersonServer(s,new(server))
	s.Serve(listenPort)
}
