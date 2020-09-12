# Comandos
```golang
// Ejemplo de comando

//mkdisk -> Sirve para crear discos
mkdisk -name->Disco1 -path->/home/archivos/discos/ -size->100 -unit->k

//rmdisk -> Sirve para eliminar discos
rmdisk -path->/home/archivos/discos/Disco1.dsk

//fdisk -> Sirve para crear , eliminar , editar particiones
//Crear Particiones
fdisk -path->/home/archivos/discos/Disco1.dsk -unit->m -size->2 -type->P -name->Part1
//Eliminar Particiones
fdisk -delete->fast -name->Part1 -path->/home/archivos/discos/Disco1.dsk
fdisk -delete->full -name->Part1 -path->/home/archivos/discos/Disco1.dsk
```

