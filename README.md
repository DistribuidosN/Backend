# 🌉 API Gateway & Middleware - Image Processing System

Este repositorio contiene el código fuente del **API Gateway (BFF - Backend for Frontend)** desarrollado en Go. Este microservicio actúa como el intermediario principal entre las aplicaciones cliente (ej. aplicación móvil en Flutter) y el Orquestador Central del sistema distribuido.

Su responsabilidad arquitectónica principal es aislar a los clientes de la complejidad de la infraestructura interna, exponiendo una interfaz RESTful moderna y traduciendo estas peticiones al protocolo SOAP requerido por el servidor central (Java), garantizando un bajo acoplamiento y alta cohesión.

## 🚀 Características Principales

* **Traducción de Protocolos:** Recibe peticiones HTTP/REST con cargas útiles en JSON y las empaqueta dinámicamente en *Envelopes* XML para su transmisión vía SOAP.
* **Procesamiento por Lotes (Batch):** Optimización de red mediante endpoints dedicados para la ingesta de múltiples imágenes en una sola petición (`UploadBatch`).
* **Soporte de Transformaciones:** Enrutamiento de parámetros para operaciones de procesamiento intensivo (ej. `GRAYSCALE`, `BLUR`, `SHARPEN`, `RESIZE`, `ROTATE`).
* **Middleware de Seguridad y Trazabilidad:** Implementación de *Recovery* nativo para prevenir caídas por *panics* y *Logger* para auditoría de tráfico.

## 🧩 Arquitectura (Flujo de Datos)

El flujo de vida de una petición sigue una estructura lineal y estandarizada:

```text
[Cliente Móvil / Web] 
       │
       ▼ (REST / JSON)
┌─────────────────────────────────┐
│       Go API Gateway            │
│  ├─ routes/   (Enrutamiento)    │
│  ├─ handlers/ (Capa HTTP/Gin)   │
│  ├─ services/ (Lógica de Negocio│
│  └─ clients/  (Cliente SOAP)    │
└─────────────────────────────────┘
       │
       ▼ (SOAP / XML)
[Servidor Central Orquestador] -> (gRPC) -> [Nodos de Procesamiento]