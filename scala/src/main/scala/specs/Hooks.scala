package simpleglobs.specs
import scala.concurrent.{ Future, ExecutionContext }
import simpleglobs.{ AkkaEnv }

trait Hooks[T] { this: AkkaEnv =>
  trait HooksSpec {
    sealed trait ChangeType
    case object Create extends ChangeType
    case object Update extends ChangeType
    case object Delete extends ChangeType

    def call(endpoint: String, data: T, changeType: ChangeType): Future[Unit]

    def callEach(endpoints: List[String], data: T, changeType: ChangeType): Future[Unit] =
      endpoints.foldLeft(Future(())) { (ops, endpoint) =>
        for {
          () <- ops
          () <- call(endpoint, data, changeType)
        } yield ()
      }
  }

  val hooks: HooksSpec
}
