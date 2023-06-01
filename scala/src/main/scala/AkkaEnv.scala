package simpleglobs
import akka.actor.typed._
import scala.concurrent.ExecutionContext

trait AkkaEnv {
  implicit val system: ActorSystem[_]
  implicit val executionContext: ExecutionContext
}
