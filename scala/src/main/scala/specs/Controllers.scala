package simpleglobs.specs
import scala.concurrent.Future

trait Controllers[T] {
  trait ControllersSpec {
    def createItem(item: T): Future[T]
    def getItem(id: String): Future[Option[T]]
    def getAllItems(): Future[List[T]]
    def putItem(id: String, item: T): Future[T]
    def deleteItem(id: String): Future[Option[T]]
  }

  val controllers: ControllersSpec
}
