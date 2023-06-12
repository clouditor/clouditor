# Demo catalog

The folder 'catalog' contains a Demo Catalog for testing purposes and can be found [here](./demo_catalog.json). The flag **all in scope** is disabled fot the Demo Catalog and the **assurance levels** *low*, *medium* and *high* are used exemplarily.

All catalogs may consist of controls, sub-controls and metrics.

The flag **all in scope** can be enabled, which is necessary for catalogs that do not allow to select the scope of the controls and thus all controls are automatically in scope. 

The controls of the catalog may contain different **assurance level** names, e.g., *low*, *medium* and *high*. These are used to decide which controls are used for the selected assurance level.

**Note:** *all in scope* and the *assurance levels* are not mutually exclusive. It is possible to select the assurance level *low* in a catalog with *all in scope* enabled. That means that all controls with assurance level *low* must be used. 
