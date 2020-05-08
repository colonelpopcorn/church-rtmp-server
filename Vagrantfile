Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/bionic64"
  config.vm.network "private_network", ip: "10.0.1.15"
  provisioner = Vagrant::Util::Platform.windows? ? :guest_ansible : :ansible
  config.vm.provision provisioner do |ansible|
    ansible.limit = "all"
    ansible.become = true
    ansible.playbook = "bootstrap.yml"
  end
end
