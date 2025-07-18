import React from 'react';
import explorerIcon from './imgs/explorador.png';

export const Menu = ({ onLogout, goToFileExplorer }) => {
    return (
        <div className="menu-container">
            <h1>Menú Principal</h1>
            <div className="icon-group">
                {/* Explorador de Archivos */}
                <div className="icon" onClick={goToFileExplorer}>
                    <img src={explorerIcon} alt="Explorador de Archivos" className="menu-icon" />
                    <p>Explorador</p>
                </div>
            </div>
            {/* Cerrar Sesión */}
            <div className="logout-container">
                <button className="logout-button" onClick={onLogout}>Cerrar Sesión</button>
            </div>
        </div>
    );
};