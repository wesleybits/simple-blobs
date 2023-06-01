package simpleglobs.defaults
import simpleglobs.specs
import simpleglobs.{ Protocol, AkkaEnv }
import scala.concurrent._

trait DefaultRepo extends specs.Repo[Protocol.Data] with AkkaEnv {
  import Protocol._
  import
    akka.actor.typed._,
    akka.actor.typed.scaladsl._

  object Repo {
    sealed trait Operations
    case class GetAll(replyTo: ActorRef[Reply]) extends Operations
    case class Get(id: String, replyTo: ActorRef[Reply]) extends Operations
    case class Put(id: String, data: Data, replyTo: ActorRef[Reply]) extends Operations
    case class Delete(id: String, replyTo: ActorRef[Reply]) extends Operations

    sealed trait Reply
    case object Ok extends Reply
    case class FoundOne(data: Option[Data]) extends Reply
    case class FoundList(data: List[Data]) extends Reply

    def apply(): Behavior[Operations] = {
      var store: Map[String, Data] = Map.empty
      Behaviors.receive { (context, message) =>
        message match {
          case GetAll(sender) => sender ! FoundList(store.values.toList)
          case Get(id, sender) => sender ! FoundOne(store.get(id))
          case Put(id, data, sender) => store += (id -> data); sender ! Ok
          case Delete(id, sender) => store -= id; sender ! Ok
        }
        Behaviors.same
      }
    }
  }

  implicit val system = ActorSystem(Repo(), "repository")
  implicit val executionContext = system.executionContext

  lazy val repo = new RepoSpec {
    import akka.actor.typed.scaladsl.AskPattern._
    import akka.util.Timeout
    import scala.concurrent.duration._

    private implicit val repository: ActorRef[Repo.Operations] = system
    private implicit val timeout = Timeout(1 seconds)
    import Repo._

    def get(id: String): Future[Option[Data]] =
      system.ask((ref:ActorRef[Reply]) => Get(id, ref)).map {
        case FoundOne(data) => data
        case _ => None
      }

    def getAll(): Future[List[Data]] =
      system.ask((ref:ActorRef[Reply]) => GetAll(ref)).map {
        case FoundList(itms) => itms
        case _ => Nil
      }

    def put(id: String, data: Data): Future[Unit] =
      system.ask((ref:ActorRef[Reply]) => Put(id, data, ref)).map { _ => () }

    def delete(id: String): Future[Unit] =
      system.ask((ref:ActorRef[Reply]) => Delete(id, ref)).map { _ => () }
  }
}

