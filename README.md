# corednsnalejplugin

## Name

*corednsnalejplugin* - CoreDNS Plugin for Nalej.

## Enabling debug mode

To enable debug mode on the plugin, configure that flag into the pluging config file:

```
    .:53 {
      corednsnalejplugin {
        systemModelAddress system-model.__NPH_NAMESPACE:8800
        debug
      }
    }
```