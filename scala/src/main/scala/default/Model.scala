package simpleglobs.defaults
import simpleglobs.specs
import simpleglobs.{ Protocol, AkkaEnv }

trait DefaultModel extends specs.Model[Protocol.Data] {
  this: specs.Repo[Protocol.Data]
      with AkkaEnv
    =>

  import Protocol._
  lazy val model = new ModelSpec {
    import scala.concurrent._

    private def genid() = java.util.UUID.randomUUID().toString

    def exists(id: String): Future[Boolean] =
      repo.get(id).map(_.isDefined)

    def create(data: Data): Future[Data] = {
      val id = genid()
      val dataWithId = data.copy(id = Some(id))
      for {
        () <- repo.put(id, dataWithId)
      } yield {
        dataWithId
      }
    }

    def update(id: String, data: Data): Future[Data] =
      for {
        dataWithId <- Future(data.copy(id = Some(id)))
        () <- repo.delete(id)
        () <- repo.put(id, dataWithId)
      } yield {
        dataWithId
      }

    def delete(id: String): Future[Option[Data]] =
      for {
        existing <- repo.get(id)
        () <- repo.delete(id)
      } yield {
        existing
      }

    def get(id: String): Future[Option[Data]] = repo.get(id)

    def getAll(): Future[List[Data]] = repo.getAll()
  }
}
