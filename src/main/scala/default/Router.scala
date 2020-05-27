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

    import akka.http.scaladsl.server._, Directives._
    import HttpMethods._
    import de.heikoseeberger.akkahttpcirce.FailFastCirceSupport._
    def routes: Route = concat(
      path("items") {
        concat(
          get {
            completeOrRecoverWith(controllers.getAllItems()) { exn =>
              printExn(exn)
              failWith(exn)
            }
          },
          post {
            entity(as[Data]) { data =>
              completeOrRecoverWith(Future(()).flatMap(_ => controllers.createItem(data))) { exn =>
                printExn(exn)
                failWith(exn)
              }
            }
          }
        )
      },
      path("item" / Segment) { id =>
        concat(
          get {
            completeOrRecoverWith(
              controllers.getItem(id)
                .map {
                  case Some(item: Data) => item
                  case None => throw new RuntimeException(s"Not Found: /item/$id")
                }
            ){ exn =>
              printExn(exn)
              failWith(exn)
            }
          },
          put {
            entity(as[Data]) { data =>
              complete {
                controllers.putItem(id, data)
              }
            }
          },
          delete {
            completeOrRecoverWith(
              controllers.deleteItem(id)
                .map {
                  case Some(item: Data) => item
                  case None => throw new RuntimeException(s"Not Found: /item/$id")
                }
            ){ exn =>
              printExn(exn)
              failWith(exn)
            }
          }
        )
      }
    )
  }
}
