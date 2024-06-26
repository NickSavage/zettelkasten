import React, { useState, useEffect } from "react";
import { getAllFiles } from "../api/files";
import { sortCards } from "../utils";
import { FileRenameModal } from "../components/FileRenameModal";
import { FileListItem } from "../components/FileListItem";

import { useNavigate } from "react-router-dom";

import { File } from "../models/File";

export function FileVault() {
  const [files, setFiles] = useState<File[]>([]);
  const [isRenameModalOpen, setIsRenameModalOpen] = useState(false);
  const [fileToRename, setFileToRename] = useState<File | null>(null);
  const navigate = useNavigate();

  function handleCardClick(id: number) {
    navigate(`/app/card/${id}`);
  }
  const openRenameModal = (file: File) => {
    setFileToRename(file);
    setIsRenameModalOpen(true);
  };

  function onDelete(file_id: number) {
    setFiles(files.filter((file) => file.id !== file_id));
  }

  function onRename(fileId: number, updatedFile: File) {
    setFiles((prevFiles) =>
      prevFiles.map((f) => (f.id === updatedFile.id ? updatedFile : f)),
    );
    setIsRenameModalOpen(false);
  }

  useEffect(() => {
    document.title = "Zettelgarden - Files";
    getAllFiles().then((data) => setFiles(sortCards(data, "sortNewOld")));
  }, []);
  return (
    <>
      <FileRenameModal
        isOpen={isRenameModalOpen}
        onClose={() => setIsRenameModalOpen(false)}
        onRename={onRename}
        file={fileToRename}
      />
      <h3>File Vault</h3>
      {files && files.length > 0 ? (
        <ul>
          {files.map((file, index) => (
            <FileListItem
              key={file.id} // Assuming each file has a unique `id` property
              file={file}
              onDelete={onDelete}
              handleViewCard={handleCardClick}
              openRenameModal={openRenameModal}
            />
          ))}
        </ul>
      ) : (
        <p>No files to display.</p> // Custom message when files array is empty
      )}
    </>
  );
}
