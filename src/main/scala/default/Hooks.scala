package simpleglobs.defaults
import simpleglobs.specs
import simpleglobs.{ Protocol, AkkaEnv }

trait DefaultHooks extends specs.Hooks[Protocol.Data] { this: AkkaEnv =>
  import Protocol._

  lazy val hooks = new HooksSpec {
    import scala.concurrent._
    import akka.http.scaladsl.Http
    import akka.http.scaladsl.model._
    import scala.util.{ Try, Success, Failure }
    private implicit val classic = system.classicSystem
    private implicit val mat = akka.stream.Materializer(classic)

    import io.circe.{ Encoder, Json }
    implicit val encodeChangeType = new Encoder[ChangeType] {
      final def apply(change: ChangeType): Json = change match {
        case Update => Json.fromString("Update")
        case Create => Json.fromString("Create")
        case Delete => Json.fromString("Delete")
      }
    }

    def call(endpoint: String, data: Data, changeType: ChangeType): Future[Unit] = {
      import io.circe.generic.auto._, io.circe.syntax._

      case class PostBody(obj: Data, change: ChangeType)
      val body = PostBody(obj = data, change = changeType).asJson.toString
      val entity = HttpEntity(ContentTypes.`application/json`, body)
      val worker = Http()
        .singleRequest(
          HttpRequest(
            uri = Uri(endpoint),
            method = HttpMethods.POST,
            entity = entity))
        .map{resp => resp.discardEntityBytes(); println(s"done: $endpoint")}

      worker.onComplete {
        case Success(_) => ()
        case Failure(err) =>
          println(err.getMessage)
          err.printStackTrace()
      }

      worker
    }
  }
}
