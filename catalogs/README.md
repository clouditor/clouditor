# Test catalogs

The following test catalogs are available for testing purposes and can be found in `./service/orchestrator/*`.
- [Example Catalog 1](./example_catalog_1.json)
- [Example Catalog 2](./example_catalog_2.json)
- [Example Catalog 3](./example_catalog_3.json)
- [Example Catalog 4](./example_catalog_4.json)

All catalogs may consist of controls, sub-controls and metrics.

The flag **all in scope** can be enabled, which is necessary for catalogs that do not allow to select the scope of the controls and thus all controls are automatically in scope. 

The controls of the catalog may contain different **assurance level** names, e.g., *low*, *medium* and *high*. These are used to decide which controls are used for the selected assurance level.

**Note:** *all in scope* and the *assurance levels* are not mutually exclusive. It is possible to select the assurance level *low* in a catalog with *all in scope* enabled. That means that all controls with assurance level *low* must be used. 

| Type | Catalog 1 | Catalog 2 | Catalog 3 | Catalog 4 |
|--- | --- | --- | --- | --- |
|Categories| 3| 3 | 3 | 3 | 3|
|Controls| yes| yes | yes | yes |
|Sub-Controls| yes| yes | yes | yes |
|Metrics| yes| yes | yes | yes |
|Assurance Levels| low <br /> medium <br /> high | basic <br /> substantial <br /> high | basic <br /> substantial <br /> high | no|
|All in scope| enabled|  enabled | disabled | disabled | 
