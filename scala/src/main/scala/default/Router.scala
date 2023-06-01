package simpleglobs.defaults
import simpleglobs.specs
import simpleglobs.{ AkkaEnv, Protocol }

trait DefaultRouter extends specs.Router {
  this: specs.Controllers[Protocol.Data] with AkkaEnv =>

  import Protocol._

  lazy val router = new RouterSpec {
    import scala.util.{ Try, Success, Failure }
    import scala.concurrent._
    import
      io.circe._,
      io.circe.generic.auto._,
      io.circe.syntax._,
      io.circe.parser._

    import
      akka.http.scaladsl._,
      akka.http.scaladsl.model._

    case class RouteError(error: String, message: String)
    private def error(message: String, exn: Throwable): RouteError =
      RouteError(message, exn.getMessage)

    private def printExn(exn: Throwable): Unit = {
      println(exn.getMessage)
      exn.printStackTrace
    }

    import akka.http.scaladsl.server._, Directives._, directives.{CompleteOrRecoverWithMagnet}
    import HttpMethods._
    import de.heikoseeberger.akkahttpcirce.FailFastCirceSupport._

    private def completeWithStdHandling(process: => CompleteOrRecoverWithMagnet) =
      completeOrRecoverWith(process){ exn =>
        printExn(exn)
        failWith(exn)
      }

    def routes: Route = concat(
      path("items") {
        concat(
          get {
            completeWithStdHandling(controllers.getAllItems())
          },
          post {
            entity(as[Data]) { data =>
              completeWithStdHandling(controllers.createItem(data))
            }
          }
        )
      },
      path("item" / Segment) { id =>
        concat(
          get {
            completeWithStdHandling(
              controllers.getItem(id)
                .map {
                  case Some(item: Data) => item
                  case None => throw new RuntimeException(s"Not Found: /item/$id")
                }
            )
          },
          put {
            entity(as[Data]) { data =>
              complete {
                controllers.putItem(id, data)
              }
            }
          },
          delete {
            completeWithStdHandling(
              controllers.deleteItem(id)
                .map {
                  case Some(item: Data) => item
                  case None => throw new RuntimeException(s"Not Found: /item/$id")
                }
            )
          }
        )
      }
    )
  }
}
