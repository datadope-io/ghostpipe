## Datasets por hacer

### Relaciones inesperadas
Relaciones topológicas cercanas por "culpa" del DNS.

Hacer un dataset como el "3" pero todo conectado al DNS.

Una posible solución, es que la distancia entre nodos pueda verse afectado por pesos en los edge y/o que el grafo sea direccional.


### Mucho ruido y pocas nueces
Meter mucho mucho ruido y tirar los servicios muy poco.

La idea es que la distancia de correlación temporal va a ser tan pequeña que no va a aportar y solo basándose en topología y etiquetas no va a encontrar nada.

Porque además, las alarmas de los backend conectados a la db van a llamarse distintas.


### Primos lejanos
Cuando se cae un servicio, se cae algo conectado a bastantes saltos de distancia.

Interconectar los ruidos con algunos saltos.

Ejemplo, se cae un bd, y donde vemos el error es en el balanceador (balanceador->frontend->backend->db)


### No me acuerdo ni de mi nombre
Intentar simular que se cae el server DNS y entonces muchos servicios se ven afectados.

Relación de topología con el DNS.


### Recien nacido
Intentamos simular que sucedería si vemos por primera vez una caída de un servicio y a que implica.

No tenemos histórico. Parecido a "mucho ruido y pocas nueces".

Por ejemplo, acaban de instalar un nuevo par de máquinas relacionadas, se cae una y la otra se ve afectada, pero este caso
nunca lo habíamos visto porque no existía.
