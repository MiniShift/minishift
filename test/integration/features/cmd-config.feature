@cmd-config @core
Feature: Minishift config subcommands
Commands `minishift config [sub-command]` are used for storing
user defined options which changes default behaviour of Minishift.

  @minishift-only @quick
  Scenario: Config is empty
    Given Minishift has state "Does Not Exist"
     Then executing "minishift config view" succeeds
      And stdout should be empty

  @quick
  Scenario Outline: Setting values in wrong range or format
  Minishift config set should not set a value, if this value is of wrong format or in wrong range.
     When executing "minishift config set <property> <value>" fails
     Then JSON config file ".minishift/config/config.json" does not contain key "<property>" with value matching "<value>"
      And stdout of command "minishift config get <property>" is equal to "<nil>"

  Examples: Wrong values for unit based keys
    | property        | value     |
    | disk-size       | 0         |
    | disk-size       | 25000.78  |
    | disk-size       | much more |
    | disk-size       | [2g 5g]   |
    | server-loglevel | 0         |
    | server-loglevel | 2500,5    |
    | server-loglevel | less      |

  Examples: Wrong values for boolean keys
    | property            | value      |
    | skip-registration   | TRuE       |
    | skip-registry-check | positive   |
    | skip-registration   | yes        |
    | skip-registry-check | fAlse      |
    | skip-registration   | -- -1      |
    | skip-registry-check | 11         |

  Examples: Wrong values for integer keys
    | property           | value      |
    | cpus               | 0          |
    | cpus               | more       |
    | cpus               | -1         |
    | cpus               | 2.5        |
    | cpus               | 2,5        |
    | server-loglevel    | 0 0        |
    | server-loglevel    | -- -1      |
    | server-loglevel    | 2.5        |
    | server-loglevel    | 2,5        |

  Examples: Wrong values for string keys
    | property          | value               |
    | password          | Mr. Three Arguments | 
    | username          | --needsquotes       |
    | no-proxy          | one two             |
    | openshift-version | two three           | 
    | public-hostname   | three four          |
    | routing-suffix    | four five           |

  Examples: Wrong values for stringSlice keys
    | property          | value               |
    | addon-env         | badly separated     |
    | insecure-registry | two, three          |
    | registry-mirror   | four five six seven |

  Examples: Wrong values for keys with extra acceptance rules
    | property        | value                          |
    | http-proxy      | bad://protocol.info            | 
    | https-proxy     | http://very.inappropriate:port |
    | host-only-cidr  | 192.168.0.1/89                 |
    | host-only-cidr  | 333.168.0.1/16                 |
    | host-only-cidr  | 192.168.1.1                    |
    | host-only-cidr  | notacidr                       |

  @quick
  Scenario Outline: Setting and unsetting with correct values
      When executing "minishift config set <property> "<value>"" succeeds
      Then JSON config file ".minishift/config/config.json" contains key "<property>" with value matching "<expected>"
       And stdout of command "minishift config get <property>" is equal to "<expected>"
       And stdout of command "minishift config view --format {{.ConfigKey}}:{{.ConfigValue}}" contains "<property>:<expected>"
      When executing "minishift config unset <property>" succeeds
      Then stdout of command "minishift config get <property>" is equal to "<nil>"
       And JSON config file ".minishift/config/config.json" does not have key "<property>"

  Examples: Correct values for unit based keys
    | property  | value   | expected |
    | disk-size | 24998   | 24998    |
    | disk-size | 24999m  | 24999m   |
    | disk-size | 25000M  | 25000M   |
    | disk-size | 25001mb | 25001mb  |
    | disk-size | 25002MB | 25002MB  |
    | disk-size | 28g     | 28g      |
    | disk-size | 29G     | 29G      |
    | disk-size | 30gb    | 30gb     |
    | disk-size | 31GB    | 31GB     |
    | memory    | 2998    | 2998     |
    | memory    | 2999m   | 2999m    |
    | memory    | 3000M   | 3000M    |
    | memory    | 3001mb  | 3001mb   |
    | memory    | 3002MB  | 3002MB   |
    | memory    | 2g      | 2g       |
    | memory    | 3G      | 3G       |
    | memory    | 4gb     | 4gb      |
    | memory    | 5GB     | 5GB      |

   Examples: Correct values for boolean keys
    | property            | value      | expected |
    | skip-registration   | true       | true     |
    | skip-registry-check | True       | true     |
    | skip-registration   | TRUE       | true     |
    | skip-registry-check | t          | true     |
    | skip-registration   | T          | true     |
    | skip-registry-check | 1          | true     |
    | skip-registration   | false      | false    |
    | skip-registry-check | False      | false    |
    | skip-registration   | FALSE      | false    |
    | skip-registry-check | f          | false    |
    | skip-registration   | F          | false    |
    | skip-registry-check | 0          | false    |
    | skip-registration   | f          | false    |
    | skip-registry-check | true       | true     |
    | skip-registration   | false      | false    |
    | skip-registry-check | True       | true     |

  Examples: Correct values for integer keys
    | property           | value      | expected |
    | cpus               | 1          | 1        |
    | cpus               | 2          | 2        |
    | cpus               | 4          | 4        |
    | server-loglevel    | 1          | 1        |
    | server-loglevel    | 5          | 5        |
    | server-loglevel    | 1111       | 1111     |

  Examples: Correct values for string keys
    | property          | value                       | expected                    |
    | password          | weakpassword                | weakpassword                |
    | password          | 2nd-And_better?             | 2nd-And_better?             |
    | username          | John MiniSnow               | John MiniSnow               |
    | username          | i.like.dots                 | i.like.dots                 |
    | no-proxy          | strange""string             | strange""string             |
    | public-hostname   | 1234567.89                  | 1234567.89                  |
    | routing-suffix    | 1234567,89                  | 1234567,89                  |

  Examples: Correct values for stringSlice keys
    | property          | value               | expected         |
    | addon-env         | one                 | [one]            |
    | docker-env        | two,three,four      | [two three four] |

  Examples: Correct values for keys with extra acceptance rules
    | property        | value                                                                                           | expected                                                                                        |
    | http-proxy      | http://proxy.io                                                                                 | http://proxy.io                                                                                 |
    | http-proxy      | http://proxy.net:3128                                                                           | http://proxy.net:3128                                                                           |
    | http-proxy      | http://me@proxy.info:91                                                                         | http://me@proxy.info:91                                                                         |
    | http-proxy      | http://me:pass@proxy.com                                                                        | http://me:pass@proxy.com                                                                        |
    | http-proxy      | http://me:pass@proxy.com:4444                                                                   | http://me:pass@proxy.com:4444                                                                   |
    | https-proxy     | https://proxy.io                                                                                | https://proxy.io                                                                                |
    | https-proxy     | https://proxy.net:3128                                                                          | https://proxy.net:3128                                                                          |
    | https-proxy     | https://me@proxy.info:91                                                                        | https://me@proxy.info:91                                                                        |
    | https-proxy     | https://me:pass@proxy.com                                                                       | https://me:pass@proxy.com                                                                       |
    | https-proxy     | https://me:pass@proxy.com:4444                                                                  | https://me:pass@proxy.com:4444                                                                  |
    | host-only-cidr  | 192.168.0.1/0                                                                                   | 192.168.0.1/0                                                                                   |
    | host-only-cidr  | 192.168.0.1/16                                                                                  | 192.168.0.1/16                                                                                  |

  @minishift-only @quick
  Scenario Outline: Setting and unsetting values for iso-url key
     When executing "minishift config set <property> "<value>"" succeeds
     Then JSON config file ".minishift/config/config.json" contains key "<property>" with value matching "<expected>"
      And stdout of command "minishift config get <property>" is equal to "<expected>"
      And stdout of command "minishift config view --format {{.ConfigKey}}:{{.ConfigValue}}" contains "<property>:<expected>"
     When executing "minishift config unset <property>" succeeds
     Then stdout of command "minishift config get <property>" is equal to "<nil>"
      And JSON config file ".minishift/config/config.json" does not have key "<property>"

  Examples: Correct values for iso-url keys
    | property        | value                                                                                            | expected                                                                                        |
    | iso-url         | https://github.com/minishift/minishift-b2d-iso/releases/download/v1.1.0/minishift-b2d.iso        | https://github.com/minishift/minishift-b2d-iso/releases/download/v1.1.0/minishift-b2d.iso       |
    | iso-url         | http://github.com/minishift/minishift-centos-iso/releases/download/v1.1.0/minishift-centos7.iso  | http://github.com/minishift/minishift-centos-iso/releases/download/v1.1.0/minishift-centos7.iso |
    | iso-url         | centos                                                                                           | centos                                                                                          |
       
  @quick
  Scenario: Setting iso-url key to a non existing path
     When executing "minishift config set iso-url file:///home/Downloads/my_handmade_centos.iso" fails
     Then stderr should contain
      """
      'file:///home/Downloads/my_handmade_centos.iso' file is not present
      """
	
  @quick
  Scenario: Unsetting non-existing key
     When executing "minishift config unset i-do-not-exist" succeeds
     Then exitcode should equal "0"

  @quick
  Scenario: Getting non-existing key
     When executing "minishift config get does-not-exist"
     Then stdout should contain "<nil>"

  @quick
  Scenario Outline: Setting values, getting values and keeping them
  Setting values, not unsetting them so they will be used on next Minishift start.
  Not every key possible is being tested only those which are less complicated,
  for example the http-proxy key is being tested in separate feature file.
     When executing "minishift config set <property> "<value>"" succeeds
     Then stdout of command "minishift config get <property>" is equal to "<expected>"

  Examples: Values to be used on next Minishift start
    | property          | value              | expected             |
    | memory            | 3500               | 3500                 |
    | disk-size         | 25g                | 25g                  |
    | docker-env        | FOO=BAR,hello=hi   | [FOO=BAR hello=hi]   |
    | docker-opt        | dns=8.8.8.8        | [dns=8.8.8.8]        |

  @quick
  Scenario Outline: Globally Setting values, getting values and keeping them
  Setting values, not unsetting them so they will be used on next Minishift start.
  Not every key possible is being tested only those which are less complicated,
  for example the http-proxy key is being tested in separate feature file.
    When executing "minishift config set --global <property> "<value>"" succeeds
    Then stdout of command "minishift config get --global <property>" is equal to "<expected>"

    Examples: Values to be used on next Minishift start
      | property          | value              | expected             |
      | memory            | 4200               | 4200                 |
      | disk-size         | 40g                | 40g                  |
      | cpus              | 3                  | 3                    |
      | docker-env        | HI=BYE,hello=hi    | [HI=BYE hello=hi]    |
      | docker-opt        | dns=1.1.1.1        | [dns=1.1.1.1]        |
      | insecure-registry | test-registry:5000 | [test-registry:5000] |

  Scenario: Minishift informs about starting with correct setup of memory, disk and CPU
  Note: Minishift rounds the values for the report to make it more readable.
        However original non-rounded values are used for the startup.
    Given Minishift has state "Does Not Exist"
      And image caching is disabled
     When executing "minishift start" succeeds
     Then stdout should match "Memory\s*:\s*3 GB"
     Then stdout should match "Disk size\s*:\s*25 GB"
     Then stdout should match "vCPUs\s*:\s*3"
      And Minishift should have state "Running"
      And Minishift VM should run with "3" vCPUs
      And Minishift VM should run within "25" to "27" GB of disk size

  Scenario: Checking that docker-env value was applied
     When printing Docker daemon configuration to stdout
     Then stdout should contain "FOO=BAR"
      And stdout should contain "hello=hi"

  Scenario: Checking that docker-opt value was applied
     When printing Docker daemon configuration to stdout
     Then stdout should contain "--dns=8.8.8.8"

  Scenario: Checking that docker-opt value was applied
     When executing "minishift ssh -- docker info"
     Then stdout should contain "test-registry:5000"

  Scenario: Deleting Minishift instance
     When executing "minishift delete --force" succeeds
     Then Minishift should have state "Does Not Exist"

  Scenario: Minishift should preserve start flags when started with non-defaults parameters.
    Given Minishift has state "Does Not Exist"
      And image caching is disabled
     When executing "minishift start --memory 5000 --disk-size 30g --cpus 2 --docker-env FOO=BAR --docker-opt dns=8.8.8.8 --insecure-registry foo.bar:5000" succeeds

  Scenario Outline: Check the config flag values
     When Minishift should have state "Running"
     Then stdout of command "minishift config get <property>" is equal to "<value>"

      Examples: Correct value show in the config
        | property          | value              |
        | memory            | 5000               |
        | disk-size         | 30g                |
        | cpus              | 2                  |
        | docker-env        | [FOO=BAR]          |
        | docker-opt        | [dns=8.8.8.8]      |
        | insecure-registry | [foo.bar:5000]     |

  Scenario: Deleting Minishift instance
    When executing "minishift delete --force" succeeds
    Then Minishift should have state "Does Not Exist"
