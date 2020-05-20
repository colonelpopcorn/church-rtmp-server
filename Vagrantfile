Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/bionic64"
  config.vm.network "private_network", ip: "10.0.1.15"
  isWindows = Vagrant::Util::Platform.windows?
  provisioner = isWindows ? :ansible_local : :ansible
  config.vm.provision provisioner do |ansible|
    ansible.limit = "all"
    ansible.sudo = true
    ansible.playbook = "bootstrap.yml"
  end
end
