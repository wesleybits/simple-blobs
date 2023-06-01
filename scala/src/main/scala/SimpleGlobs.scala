package simpleglobs

import simpleglobs.specs._
import simpleglobs.defaults._
import akka.actor._
import akka.stream._
import akka.http.scaladsl._

class SimpleGlobs
    extends AkkaEnv
    with DefaultModel
    with DefaultRepo
    with DefaultHooks
    with DefaultControllers
    with DefaultRouter

object SimpleGlobs {
  def main(args: Array[String]): Unit = {

    val app = new SimpleGlobs()
    import app.system
    import system.executionContext
    implicit val classic = app.system.classicSystem
    implicit val mat = akka.stream.Materializer(classic)

    val serverBinding = Http().bindAndHandle(app.router.routes, "localhost", 8080)
    println(s"SimpleGlobs serviced at: http://localhost:8080\nPress RETURN to stop.")
    scala.io.StdIn.readLine()
    serverBinding
      .flatMap(_.unbind())
      .onComplete(_ => system.terminate())
  }
}
