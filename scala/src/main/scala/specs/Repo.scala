package simpleglobs.specs

import akka.actor._
import scala.concurrent._


trait Repo[T] {
  trait RepoSpec {
    def get(id: String): Future[Option[T]]
    def getAll(): Future[List[T]]
    def put(id: String, data: T): Future[Unit]
    def delete(id: String): Future[Unit]
  }

  val repo: RepoSpec
}
