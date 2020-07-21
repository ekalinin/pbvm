pbvm
====

pbvm - Protocol Buffers Version Manager.

Install
=======

```sh
# download & unpack
$ wget https://github.com/ekalinin/pbvm/releases/download/v0.0.1/pbvm_0.0.1_linux_x86_64.tar.gz
$ tar pbvm_0.0.1_linux_x86_64.tar.gz

# install
$ sudo cp pbvm /usr/bin

# activate
$ export PATH="$PATH:$HOME/.pbvm/active/bin"
```

Usage
=====

List available versions
-----------------------

```sh
$ pbvm list-remote -n 5     
   VERSION   | PRE-RELEASE |    DATE    | INSTALLED  
-------------+-------------+------------+------------
  v4.0.0-rc1 | true        | 2020.07.15 | true       
  v3.12.3    | false       | 2020.06.03 | true       
  v3.12.2    | false       | 2020.05.26 | false      
  v3.12.1    | false       | 2020.05.20 | false      
  v3.12.0    | false       | 2020.05.15 | true 
```

Install (switch) to version
----------------------------

```sh
$ pbvm install v4.0.0-rc1
$ protoc --version
libprotoc 4.0.0

$ pbvm install v3.12.3
$ protoc --version
libprotoc 3.12.3

# will just switch active version (without downloading)
$ pbvm install v4.0.0-rc1
$ protoc --version
libprotoc 3.12.3
```

List local versions
-------------------

```sh
$ pbvm ls                                         
   VERSION   | INSTALL DATE | ACTIVE  
-------------+--------------+---------
  v4.0.0-rc1 | 2020.07.20   | false   
  v3.12.3    | 2020.07.20   | true    
  v3.12.0    | 2020.07.21   | false 
```

Run with a version
------------------

```sh
$ protoc --version
libprotoc 3.12.3

$ pbvm run "protoc --version" --version v4.0.0-rc1        
libprotoc 4.0.0

$ protoc --version
libprotoc 3.12.3
```

Auto completion
---------------

```sh
# see instructions below
$ pbvm completion -h
```