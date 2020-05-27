package simpleglobs
import io.circe.Json

object Protocol {

  case class Data(
    hooks: List[String],
    data: Json,
    id: Option[String])
}
