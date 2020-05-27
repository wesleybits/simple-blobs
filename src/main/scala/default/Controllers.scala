package simpleglobs.defaults
import simpleglobs.{ Protocol, AkkaEnv }
import simpleglobs.specs

trait DefaultControllers extends specs.Controllers[Protocol.Data] {
  this: specs.Model[Protocol.Data]
      with specs.Hooks[Protocol.Data]
      with AkkaEnv =>

  import Protocol._

  lazy val controllers = new ControllersSpec {
    import scala.concurrent.Future

    def createItem(item: Data): Future[Data] = model.create(item).map { item =>
      val _hooksWorker = hooks.callEach(item.hooks, item, hooks.Create)
      item
    }

    def getItem(id: String): Future[Option[Data]] = model.get(id)

    def getAllItems(): Future[List[Data]] = model.getAll()

    def putItem(id: String, item: Data): Future[Data] = for {
      existing <- model.get(id)
      updated <- model.update(id, item)
    } yield {
      val allHooks =
        (existing.fold(List.empty[String])(_.hooks) ++ updated.hooks)
          .distinct
      val _hooksWorker = hooks.callEach(allHooks, updated, hooks.Update)
      updated
    }

    def deleteItem(id: String): Future[Option[Data]] = for {
      item <- model.delete(id)
    } yield {
      val allHooks = item.fold(List.empty[String])(_.hooks)
      val _hooksWorker = item.fold(Future(()))(hooks.callEach(allHooks, _, hooks.Delete))
      item
    }
  }
}
