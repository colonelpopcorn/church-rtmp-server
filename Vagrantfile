Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/bionic64"
  config.vm.network "private_network", ip: "10.0.1.15"
  isWindows = Vagrant::Util::Platform.windows?
  provisioner = isWindows ? :guest_ansible : :ansible
  config.vm.provision provisioner do |ansible|
    ansible.limit = "all"
    # TODO: Need to update plugin upstream to come to parity with Vagrant's
    # ansible interface.
    if isWindows
      ansible.sudo = true
    else
      ansible.become = true
    end
    ansible.playbook = "bootstrap.yml"
  end
end
