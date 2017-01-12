# Routes service

Service that provides BGP routes and neighboring from a set of virtual
appliances that wants to talk BGP.

As an example, imagine a set of virtualized L4 Balancers that wants to speak
BGP to a fisical Router to expose their routes. As the balancers are
dinamically created, you should not peer directly to the fisical router
without configuring it or enabling a security constraint by allowing anything
that speaks BGP in your network to peer with your fisical router.

For that, one can peer the Routes Service with a fisical router. Then connect
all virtual appliances to the service that will expose neighbors and routes to
the fisical router.

To achieve that, we expose a HTTP service to be used by the dinamic
appliances. Internally, the service uses [gobgp](https://github.com/osrg/gobgp) to run a BGP deamon that peers
with fisical router.
