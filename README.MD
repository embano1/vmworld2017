# Repository for VMworld US 2017 Session [DEV2133BU]
Contains code, descriptions and ready to use Kubernetes manifests.  
Requirements:

- Kubernetes environment running (e.g. minikube and kubectl, tested with Kubernetes and kubectl v1.6.4)
- go tools (compiler, etc.) and `make` in case you want to play with the code
- `go-wrk` [https://github.com/tsliwowicz/go-wrk] benchmark if you want to follow along with the rolling update example (01_probe)

# Session abstract
**Designing an Application-Centric Infrastructure Platform: Experience from the Field [DEV2133BU]**  
Michael Gasch, Application Platform Architect, VMware  
Martijn Baecke, VMware  

In the age of digital transformation, businesses want a competitive advantage. They rely heavily on agile applications with just-in-time data intelligence. To meet these expectations, architects must provide platforms that are agile across the application layer and infrastructure, along with operational excellence. It is imperative that developers and infrastructure operations work closely together to build the necessary architecture with flexibility and robustness in mind from the beginning. In this session, Emad Benjamin and Martijn Baecke will share how to design, build, and operate a platform that can run any application; discuss best practices in the field from both application and infrastructure perspectives; and address topics like DevOps, site reliability engineering, and operational transformation.

# Play with the examples
Each sub-folder has it´s own README, with descriptions how to run the demos, get and build the code, etc.