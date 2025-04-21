# -- Vagrantfile --
Vagrant.configure("2") do |config|
  # 1) Base box and hostname
  config.vm.box      = "ubuntu/focal64"
  config.vm.hostname = "cfm-test"

  # 2) Networking
  config.vm.network "private_network", ip: "192.168.56.11"
  config.vm.network "forwarded_port", guest: 80, host: 8080

  # 3) Sync your local cfm/ folder
  config.vm.synced_folder "./cfm", "/home/vagrant/cfm"

  # 4) VM resources
  config.vm.provider "virtualbox" do |vb|
    vb.memory = "2048"
  end

  # 5) Provisioning script
  config.vm.provision "shell", inline: <<-SHELL
    set -eux

    # --- A) Install Docker CE & Compose
    apt-get update
    apt-get install -y \
      ca-certificates \
      curl \
      gnupg \
      lsb-release

    # Add Docker’s official GPG key
    install -m0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg \
      | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    chmod a+r /etc/apt/keyrings/docker.gpg

    # Set up the Docker apt repo
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
       https://download.docker.com/linux/ubuntu \
       $(lsb_release -cs) stable" \
      > /etc/apt/sources.list.d/docker.list

    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io

    # Install Docker Compose v2
    COMPOSE_VERSION=v2.1.1
    curl -fsSL "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" \
      -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose

    # --- B) Install and configure Nginx
    apt-get install -y nginx

    # --- C) Add vagrant to the docker group
    usermod -aG docker vagrant

    # --- D) Build & run your Go HTTP wrapper container
    cd /home/vagrant/cfm

    docker build -t cfmid-wrapper:local .
    docker rm -f cfmid-wrapper || true
    docker run -d \
      --name cfmid-wrapper \
      --network host \
      --restart unless-stopped \
      cfmid-wrapper:local

    # --- E) Deploy nginx.conf if present
    if [ -f /home/vagrant/cfm/nginx.conf ]; then
      install -m644 /home/vagrant/cfm/nginx.conf /etc/nginx/sites-available/cfm
      ln -sf /etc/nginx/sites-available/cfm /etc/nginx/sites-enabled/cfm
      rm -f /etc/nginx/sites-enabled/default
      nginx -s reload
    fi

  SHELL
end