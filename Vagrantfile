# -*- mode: ruby -*-
# vi: set ft=ruby :

# VM used to test xdp program
Vagrant.configure("2") do |config|
    config.vm.box = "ubuntu/jammy64"
    
    config.vm.define "xdp", primary: true do |xdp|
        xdp.vm.network "private_network", ip: "10.10.10.10"
        xdp.vm.synced_folder "./", "/home/vagrant/xdpdropper"
        xdp.vm.provision "shell", inline: <<-SHELL
            apt -y update
            apt install -y \
                golang \
                make \
                ca-certificates \
                curl \
                gnupg \
                lsb-release
            mkdir -p /etc/apt/keyrings
            curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
            echo \
                "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
                $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
            apt -y update
            apt install -y \
                docker-ce \
                docker-ce-cli \
                containerd.io \
                docker-compose-plugin
            usermod -aG docker vagrant
            curl -sLS https://get.k3sup.dev | sh
            sudo install k3sup /usr/local/bin/
        SHELL
        xdp.vm.network "forwarded_port", guest: 8080, host: 8080
    end

    config.vm.define "test", primary: true do |test|
        test.vm.network "private_network", ip: "10.10.10.20"
    end
end
