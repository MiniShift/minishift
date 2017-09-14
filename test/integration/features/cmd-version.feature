@cmd-version @command @openshift @minishift-only
Feature: Version command
Command "minishift version" shows version of minishift binary.

  Scenario: User can get information about version of Minishift
      When executing "minishift version" succeeds
      Then stdout should match "^minishift v[0-9]\.[0-9]\.[0-9]\+[a-z0-9]{7,8}\n$"
