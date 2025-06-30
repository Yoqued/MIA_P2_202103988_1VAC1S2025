import React, { useEffect, useState } from 'react';
import diskIcon from './imgs/disk.png';
import partitionIcon from './imgs/partition.png';
import folderIcon from './imgs/carpeta.png';
import fileIcon from './imgs/file.png';

export const Explorador = ({ onBackToMenu }) => {
    const [discos, setDiscos] = useState([]);
    const [particiones, setParticiones] = useState([]);
    const [files, setFiles] = useState([]);
    const [content, setContent] = useState("");
    const [selectedDisco, setSelectedDisco] = useState(null);
    const [selectedPartition, setSelectedPartition] = useState(null);
    const [currentPath, setCurrentPath] = useState("/");
    const [error, setError] = useState('');

    // Función para obtener los discos del backend
    const fetchDiscos = async () => {
        try {
            const response = await fetch('http://3.142.172.171:8080/discos', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                console.log(data)
                setDiscos(data.discos);
            } else {
                setError('Error al obtener los discos');
            }
        } catch (err) {
            setError('Error al conectar con el servidor.');
        }
    };

    // UseEffect para cargar los discos cuando el componente se monta
    useEffect(() => {
        fetchDiscos();
    }, []);

    // Función para manejar cuando se hace click en un disco
    const handleDiscoClick = async (disco) => {
        setError('');
        try {
            const response = await fetch(`http://3.142.172.171:8080/disco-select?path=${encodeURIComponent(disco.Path)}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                console.log("Información del disco seleccionada:", data);
                setParticiones(data.filter(part => part.Name && part.Name.trim() !== "\u0000"));
                setSelectedDisco(disco);
            } else {
                setError('Error al obtener los detalles del disco.');
            }
        } catch (err) {
            setError('Error al conectar con el servidor.');
        }
    };

    // Función para manejar cuando se hace click en una partición
    const handlePartitionClick = async (particion, disco) => {
        setError('');
        try {
            const response = await fetch(`http://3.142.172.171:8080/partition-select?name=${encodeURIComponent(particion.Name)}
            &start=${encodeURIComponent(particion.Start)}&path=${encodeURIComponent(disco.Path)}&fileparts=${encodeURIComponent(currentPath)}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                setFiles(data);
                setSelectedPartition(particion);
            } else {
                setError('Error al obtener los detalles del disco.');
            }
        } catch (err) {
            setError('Error al conectar con el servidor.');
        }
    };

    // Función para manejar cuando se hace click en una carpeta/archivo
    const handleFileClick = async (name, disco, particion, type) => {
        setError('');
        const newPath = currentPath + name + "/";  // Actualiza la ruta actual
        setCurrentPath(newPath);

        if (type === 'carpeta') {
            // Si es una carpeta, cargar los archivos de esa carpeta
            try {
                const response = await fetch(`http://3.142.172.171:8080/file-select?start=${encodeURIComponent(particion.Start)}
                &path=${encodeURIComponent(disco.Path)}&fileparts=${encodeURIComponent(newPath)}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });

                if (response.ok) {
                    const data = await response.json();
                    setFiles(data);  // Actualizar la lista de archivos con el nuevo contenido
                } else {
                    setError('Error al obtener los archivos.');
                }
            } catch (err) {
                setError('Error al conectar con el servidor.');
            }
        } else {
            // Si es un archivo de texto, cargar el contenido del archivo
            handleContentClick(name, disco, particion);
        }
    };

    const handleContentClick = async (name, disco, particion) => {
        setError('');
        const newPath = currentPath + name;  // Calcula la nueva ruta antes de actualizar el estado

        try {
            const response = await fetch(`http://3.142.172.171:8080/content-select?start=${encodeURIComponent(particion.Start)}
            &path=${encodeURIComponent(disco.Path)}&fileparts=${encodeURIComponent(newPath)}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                setContent(data);  // Establecer el contenido del archivo de texto
            } else {
                setError('Error al obtener los detalles del archivo.');
            }
        } catch (err) {
            setError('Error al conectar con el servidor.');
        }
    };

    const handleBackFile = async () => {
        setError('');
        if (content) {
            // Si está mostrando contenido, ocultarlo y mostrar la lista de archivos de nuevo
            setContent("");
        } 
        if (currentPath !== "/") {
            // Elimina la última carpeta del currentPath
            const newPath = currentPath.slice(0, currentPath.lastIndexOf("/", currentPath.length - 2) + 1);
            setCurrentPath(newPath);

            // Realiza la solicitud de nuevo para obtener los archivos de la carpeta superior
            try {
                const response = await fetch(`http://3.142.172.171:8080/file-select?start=${encodeURIComponent(selectedPartition.Start)}
                &path=${encodeURIComponent(selectedDisco.Path)}&fileparts=${encodeURIComponent(newPath)}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });

                if (response.ok) {
                    const data = await response.json();
                    setFiles(data);  // Actualizar la lista de archivos con el nuevo contenido
                } else {
                    setError('Error al obtener los archivos.');
                }
            } catch (err) {
                setError('Error al conectar con el servidor.');
            }
        }
    };

    const handleBackToPartitions = () => {
        setError('');
        setSelectedPartition(null);  // Limpiar la partición seleccionada
        setFiles([]);
        setContent("");               // Limpiar el contenido
        setCurrentPath("/"); 
    };

    const handleBackToDiscos = () => {
        setError('');
        setSelectedDisco(null); // Limpiar el disco seleccionado y regresar a la vista de discos
        setParticiones([]); // Limpiar la lista de particiones
    };

    return (
        <div className="explorador-container">
            <h1>Explorador de Archivos</h1>
            {error && <p className="error">{error}</p>}
            
            {/* Mostrar discos, particiones o archivos dependiendo del estado */}
            <div className='contenedor'>
                {!selectedDisco ? (
                    <>
                        <div className="discos-list">
                            {discos.length > 0 ? (
                                discos.map((disco, index) => (
                                    <div key={index} className="disco-item" onClick={() => handleDiscoClick(disco)}>
                                        <img src={diskIcon} alt="Disk Icon" />
                                        <p>{disco.Name}</p>
                                    </div>
                                ))
                            ) : (
                                <p>No hay discos disponibles</p>
                            )}
                        </div>
                        <button onClick={onBackToMenu} className="button-exit">Regresar al Menú</button> 
                    </>
                ) : !selectedPartition ? (
                    <>
                        <h2>Particiones de {selectedDisco.Name}</h2>
                        <div className="particiones-list">
                            {particiones.length > 0 ? (
                                particiones.map((particion, index) => (
                                    <div 
                                        key={index} 
                                        className="particion-item" 
                                        onClick={() => {
                                            if (particion.Status === "0") {
                                                alert("La partición no está montada. No es posible acceder a sus archivos.");
                                            } else {
                                                handlePartitionClick(particion, selectedDisco);
                                            }
                                        }}
                                    >
                                        <img src={partitionIcon} alt="Partition Icon" />
                                        <p>Nombre: {particion.Name}</p>
                                        <p>Estado: {particion.Status === "1" ? "Montada" : "No montada"}</p>
                                    </div>
                                ))
                            ) : (
                                <p>No hay particiones disponibles</p>
                            )}
                        </div>
                        <button className="button" onClick={handleBackToDiscos}>Regresar a discos</button>
                    </>
                ) : (
                    <>
                        {/* Barra de navegación */}
                        <div className="file-path-bar">
                            <p>{currentPath}</p> 
                        </div>

                        {/* Mostrar el contenido del archivo si está disponible */}
                        {content ? (
                            <div className="file-content">
                                <h3>Contenido del archivo:</h3>
                                <textarea readOnly value={content} rows="10" cols="50"></textarea>
                            </div>
                        ) : (
                            <div className="files-list">
                                {files.length > 0 ? (
                                    files.map((file, index) => (
                                        <div key={index} className="file-item">
                                            <img src={file.Type === 'carpeta' ? folderIcon : fileIcon} alt="File Icon" 
                                            onClick={() => handleFileClick(file.Name, selectedDisco, selectedPartition, file.Type)} />
                                            <p>{file.Name}</p>
                                        </div>
                                    ))
                                ) : (
                                    <p>No hay archivos disponibles</p>
                                )}
                            </div>
                        )}
                        <button className="button" onClick={handleBackFile}>Regresar a carpeta anterior</button>
                        <button className="button" onClick={handleBackToPartitions}>Regresar a particiones</button>
                    </>
                )}
            </div>
        </div>
    );
};