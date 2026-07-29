resource "docker_image" "nginx" {
  name = "nginx:latest"

}

resource "docker_container" "nginx" {
  name = "platform-nginx"

  image = docker_image.nginx.image_id

  restart = "unless-stopped"

  ports {
    internal = 80
    external = 8080
  }

  networks_advanced {
    name = docker_network.platform.name
  }
}


