Vagrant.configure("2") do |config|
  config.vm.box = "boxomatic/alpine-3.14"
  config.vm.network "private_network", ip: "10.0.1.15"
  config.vm.hostname = 'streaming-server.vagrant'
  isWindows = Vagrant::Util::Platform.windows?
  provisioner = isWindows ? :ansible_local : :ansible
  config.vm.provision provisioner do |ansible|
    ansible.limit = "all"
    ansible.become = true
    ansible.playbook = "provisioning/bootstrap.yml"
    ansible.extra_vars = {
      owner: "vagrant"
      is_vm: true
    }
    # Change and uncomment to work only on certain parts.
    # ansible.raw_arguments = ["--tags", "frontend"]
  end
end
