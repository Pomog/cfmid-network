# -- mode: ruby --
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/focal64"
  config.vm.hostname = "cfm-test"

  config.vm.network "private_network", ip: "192.168.56.11"

  config.vm.network "public_network"

  config.vm.network "forwarded_port", guest: 80, host: 8080

  config.vm.synced_folder "./cfm", "/home/vagrant/cfm"

  config.vm.provider "virtualbox" do |vb|
     vb.memory = "2048"
   end

  # 4) Provision: install Docker & Nginx, build & run container, configure proxy
  config.vm.provision "shell", inline: <<-SHELL
    sudo apt-get update
    sudo apt-get install \
        ca-certificates \
        curl \
        gnupg -y

    sudo install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    sudo chmod a+r /etc/apt/keyrings/docker.gpg
    echo \
      "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update
    sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.1.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

    chmod +x /usr/local/bin/docker-compose

    docker run wishartlab/cfmid

    SHELL
end