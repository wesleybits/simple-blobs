val akkaVersion = "2.6.5"
val circeVersion = "0.13.0"

val akkaOrg = "com.typesafe.akka"
val akkaDeps = Seq(
  akkaOrg %% "akka-actor-typed" % akkaVersion,
  akkaOrg %% "akka-stream" % akkaVersion )

val circeOrg = "io.circe"
val circeDeps = Seq(
  circeOrg %% "circe-core" % circeVersion,
  circeOrg %% "circe-generic" % circeVersion,
  circeOrg %% "circe-parser" % circeVersion )

val httpDeps = Seq(
  akkaOrg %% "akka-http" % "10.1.12",
  "de.heikoseeberger" %% "akka-http-circe" % "1.32.0" )

lazy val root = (project in file("."))
  .settings(
    name := "simpleglobs",
    scalaVersion := "2.12.7",
    libraryDependencies ++= akkaDeps ++ circeDeps ++ httpDeps,
    resourceDirectory := baseDirectory.value / "resource")
