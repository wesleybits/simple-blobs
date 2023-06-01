package simpleglobs.specs
import akka.http.scaladsl.server.Route
import scala.concurrent.Future

trait Router {
  trait RouterSpec {
    def routes: Route
  }

  val router: RouterSpec
}
