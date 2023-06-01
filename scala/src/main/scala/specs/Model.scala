package simpleglobs.specs
import scala.concurrent._

trait Model[T] {
  trait ModelSpec {
    def exists(id: String): Future[Boolean]
    def update(id: String, item: T): Future[T]
    def create(item: T): Future[T]
    def get(id: String): Future[Option[T]]
    def getAll(): Future[List[T]]
    def delete(id: String): Future[Option[T]]
  }

  val model: ModelSpec
}
