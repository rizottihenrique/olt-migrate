"use client";

import { useState } from "react";

export default function Home() {
  const [file, setFile] = useState<File | null>(null);
  const [inputMode, setInputMode] = useState<"file" | "paste">("file");
  const [pastedText, setPastedText] = useState("");
  const [isDragging, setIsDragging] = useState(false);
  const [sourceVendor, setSourceVendor] = useState("Fiberhome AN5000");
  const [destVendor, setDestVendor] = useState("Nokia");
  const [mappings, setMappings] = useState([{ sourceSlot: 1, sourcePon: 1, destSlot: 1, destPon: 1 }]);
  const [commands, setCommands] = useState("");
  const [ponCommands, setPonCommands] = useState<Record<string, string>>({});
  const [onuCommands, setOnuCommands] = useState<Record<string, string>>({});
  const [viewMode, setViewMode] = useState<"all" | "clean" | "pon" | "onu">("all");
  const [selectedPon, setSelectedPon] = useState<string>("");
  const [searchOnu, setSearchOnu] = useState<string>("");
  const [hideMarkers, setHideMarkers] = useState<boolean>(false);
  const [loading, setLoading] = useState(false);
  const [stats, setStats] = useState<{ totalONUsFound: number } | null>(null);

  const getCleanScript = (text: string) => {
    if (!text) return "";
    return text
      .split("\n")
      .filter((line) => {
        const trimmed = line.trim();
        if (trimmed.startsWith("! ONU") || trimmed.startsWith("! ONT") || trimmed.startsWith("! --- Provisionamento da ONU") || trimmed.startsWith("! --- ONU")) {
          return false;
        }
        return true;
      })
      .join("\n")
      .replace(/\n{3,}/g, "\n\n");
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
    }
  };

  const updateMapping = (index: number, field: string, value: string) => {
    const newMappings = [...mappings];
    newMappings[index] = { ...newMappings[index], [field]: parseInt(value) || 0 };
    setMappings(newMappings);
  };

  const addMapping = () => {
    setMappings([...mappings, { sourceSlot: 1, sourcePon: 1, destSlot: 1, destPon: 1 }]);
  };

  const removeMapping = (index: number) => {
    const newMappings = mappings.filter((_, i) => i !== index);
    setMappings(newMappings);
  };

  const handleMigrate = async () => {
    let targetFile = file;
    if (inputMode === "paste") {
      if (!pastedText.trim()) return alert("Cole o texto de configuração da Fiberhome no campo.");
      targetFile = new File([pastedText], "backup_colado.txt", { type: "text/plain" });
    } else if (!targetFile) {
      return alert("Selecione um arquivo de configuração da Fiberhome.");
    }
    
    setLoading(true);
    const formData = new FormData();
    formData.append("configFile", targetFile);
    formData.append("sourceVendor", sourceVendor);
    formData.append("destVendor", destVendor);
    formData.append("mappings", JSON.stringify(mappings));

    try {
      const res = await fetch("/api/migrate", {
        method: "POST",
        body: formData,
      });

      if (!res.ok) {
        throw new Error("Erro na migração: " + await res.text());
      }

      const data = await res.json();
      setCommands(data.commands || "");
      setPonCommands(data.ponCommands || {});
      setOnuCommands(data.onuCommands || {});
      setStats({ totalONUsFound: data.totalONUsFound });
      setViewMode("all");
      const pons = Object.keys(data.ponCommands || {});
      if (pons.length > 0) setSelectedPon(pons[0]);
    } catch (err: any) {
      alert(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = () => {
    if (!commands) return;
    const blob = new Blob([commands], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `script_migracao_${sourceVendor}_para_${destVendor}.txt`.toLowerCase();
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <div className="flex flex-col gap-6 lg:flex-row items-start">
      
      {/* Coluna Esquerda: Configurações */}
      <div className="w-full lg:w-5/12 flex flex-col gap-6">
        
        {/* Passo 1: Fabricantes */}
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden flex flex-col">
          <div className="px-5 py-4 border-b border-gray-100 bg-gray-50/50">
            <h2 className="text-sm font-semibold text-gray-800">1. Fabricantes de OLT</h2>
          </div>
          <div className="p-5 flex gap-4">
            <div className="flex-1">
              <label className="text-[11px] font-semibold text-gray-500 uppercase tracking-wider mb-1 block">Origem (De)</label>
              <select value={sourceVendor} onChange={(e) => setSourceVendor(e.target.value)} className="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none bg-white">
                <option value="Fiberhome AN5000">Fiberhome AN5000</option>
                <option value="Fiberhome AN6000">Fiberhome AN6000</option>
                <option value="ZTE">ZTE (Em breve)</option>
                <option value="Huawei">Huawei (Em breve)</option>
              </select>
            </div>
            <div className="flex-1">
              <label className="text-[11px] font-semibold text-gray-500 uppercase tracking-wider mb-1 block">Destino (Para)</label>
              <select value={destVendor} onChange={(e) => setDestVendor(e.target.value)} className="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none bg-white">
                <option value="Nokia">Nokia ISAM</option>
                <option value="Fiberhome AN5000">Fiberhome AN5000</option>
                <option value="Fiberhome AN6000">Fiberhome AN6000</option>
                <option value="Huawei">Huawei (Em breve)</option>
                <option value="Datacom">Datacom (Em breve)</option>
              </select>
            </div>
          </div>
        </div>

        {/* Passo 2: Arquivo ou Texto */}
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
          <div className="px-5 py-3 border-b border-gray-100 bg-gray-50/50 flex justify-between items-center">
            <h2 className="text-sm font-semibold text-gray-800">2. Backup (Arquivo ou Texto Colado)</h2>
            <div className="flex bg-gray-200/60 p-0.5 rounded-lg text-xs font-medium">
              <button
                type="button"
                onClick={() => setInputMode("file")}
                className={`px-2.5 py-1 rounded-md transition-all ${
                  inputMode === "file" ? "bg-white text-blue-600 shadow-sm font-semibold" : "text-gray-600 hover:text-gray-900"
                }`}
              >
                📂 Arquivo
              </button>
              <button
                type="button"
                onClick={() => setInputMode("paste")}
                className={`px-2.5 py-1 rounded-md transition-all ${
                  inputMode === "paste" ? "bg-white text-blue-600 shadow-sm font-semibold" : "text-gray-600 hover:text-gray-900"
                }`}
              >
                📋 Colar Texto
              </button>
            </div>
          </div>
          <div className="p-5">
            {inputMode === "file" ? (
              <label 
                onDragOver={(e) => { e.preventDefault(); setIsDragging(true); }}
                onDragLeave={(e) => { e.preventDefault(); setIsDragging(false); }}
                onDrop={(e) => {
                  e.preventDefault();
                  setIsDragging(false);
                  if (e.dataTransfer.files && e.dataTransfer.files[0]) {
                    setFile(e.dataTransfer.files[0]);
                  }
                }}
                className={`block w-full border-2 border-dashed transition-colors rounded-lg p-6 text-center cursor-pointer ${
                  isDragging ? "border-blue-500 bg-blue-50/80" : "border-gray-200 hover:border-blue-400 hover:bg-blue-50/50"
                }`}
              >
                <input type="file" onChange={handleFileChange} className="hidden" />
                {file ? (
                  <div>
                    <svg className="w-8 h-8 mx-auto mb-2 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>
                    <p className="text-sm font-medium text-gray-800">{file.name}</p>
                    <p className="text-xs text-gray-500 mt-1">{(file.size / 1024).toFixed(1)} KB</p>
                    <span className="inline-block mt-2 text-[11px] text-blue-600 font-medium underline">Clique ou arraste para trocar o arquivo</span>
                  </div>
                ) : (
                  <div>
                    <svg className="w-8 h-8 mx-auto mb-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
                    <p className="text-sm font-medium text-gray-700">Clique para escolher ou arraste o arquivo aqui</p>
                    <p className="text-xs text-gray-400 mt-1">Aceita arquivos .txt, .cfg, .dat, .log ou sem extensão</p>
                  </div>
                )}
              </label>
            ) : (
              <div>
                <textarea
                  value={pastedText}
                  onChange={(e) => setPastedText(e.target.value)}
                  placeholder="Cole aqui o texto do seu backup ou comandos da OLT (ex: ! ONU 2... ou whitelist add...)"
                  className="w-full h-36 p-3 text-xs font-mono bg-gray-50 border border-gray-300 rounded-lg focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none resize-y transition-shadow"
                />
                <div className="flex justify-between items-center mt-1.5 text-xs text-gray-500">
                  <span>{pastedText ? `${pastedText.split("\n").length} linhas` : "Nenhum texto colado"}</span>
                  {pastedText && (
                    <button type="button" onClick={() => setPastedText("")} className="text-red-500 hover:underline">
                      Limpar texto
                    </button>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Passo 3: Mapeamento */}
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden flex flex-col">
          <div className="px-5 py-4 border-b border-gray-100 bg-gray-50/50 flex justify-between items-center">
            <h2 className="text-sm font-semibold text-gray-800">3. De-Para (Portas PON)</h2>
            <button onClick={addMapping} className="text-xs font-medium text-blue-600 hover:text-blue-800">
              + Adicionar regra
            </button>
          </div>
          
          <div className="p-5 space-y-4">
            {mappings.map((m, i) => (
              <div key={i} className="flex flex-col gap-2 p-3 bg-gray-50 rounded-lg border border-gray-100">
                <div className="flex items-center gap-4">
                  <div className="flex-1">
                    <label className="text-[11px] font-semibold text-gray-500 uppercase tracking-wider mb-1 block">{sourceVendor}</label>
                    <div className="flex gap-2">
                      <div className="flex-1">
                        <span className="text-xs text-gray-400 absolute ml-2 mt-1.5">Sl</span>
                        <input type="number" min="1" value={m.sourceSlot} onChange={(e) => updateMapping(i, "sourceSlot", e.target.value)} className="w-full pl-6 pr-2 py-1.5 text-sm border border-gray-300 rounded focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-shadow bg-white" />
                      </div>
                      <div className="flex-1">
                        <span className="text-xs text-gray-400 absolute ml-2 mt-1.5">Pn</span>
                        <input type="number" min="1" value={m.sourcePon} onChange={(e) => updateMapping(i, "sourcePon", e.target.value)} className="w-full pl-7 pr-2 py-1.5 text-sm border border-gray-300 rounded focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-shadow bg-white" />
                      </div>
                    </div>
                  </div>

                  <div className="mt-5 text-gray-300">
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 5l7 7m0 0l-7 7m7-7H3" /></svg>
                  </div>

                  <div className="flex-1">
                    <label className="text-[11px] font-semibold text-blue-600 uppercase tracking-wider mb-1 block">{destVendor}</label>
                    <div className="flex gap-2">
                      <div className="flex-1">
                        <span className="text-xs text-blue-400/50 absolute ml-2 mt-1.5">Sl</span>
                        <input type="number" min="1" value={m.destSlot} onChange={(e) => updateMapping(i, "destSlot", e.target.value)} className="w-full pl-6 pr-2 py-1.5 text-sm border border-blue-200 rounded focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-shadow bg-blue-50/30" />
                      </div>
                      <div className="flex-1">
                        <span className="text-xs text-blue-400/50 absolute ml-2 mt-1.5">Pn</span>
                        <input type="number" min="1" value={m.destPon} onChange={(e) => updateMapping(i, "destPon", e.target.value)} className="w-full pl-7 pr-2 py-1.5 text-sm border border-blue-200 rounded focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-shadow bg-blue-50/30" />
                      </div>
                    </div>
                  </div>
                </div>
                {mappings.length > 1 && (
                  <button onClick={() => removeMapping(i)} className="text-[11px] text-red-500 hover:text-red-700 text-left mt-1 font-medium">
                    Remover
                  </button>
                )}
              </div>
            ))}

            <button 
              onClick={handleMigrate} 
              disabled={loading || (inputMode === "file" ? !file : !pastedText.trim())}
              className={`w-full py-2.5 rounded-lg text-sm font-medium text-white shadow-sm transition-all mt-4 flex items-center justify-center
                ${loading || (inputMode === "file" ? !file : !pastedText.trim()) ? 'bg-gray-300 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700 hover:shadow'}`}
            >
              {loading ? (
                <>
                  <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24"><circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle><path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                  Processando...
                </>
              ) : (
                "Gerar Scripts de Migração"
              )}
            </button>
          </div>
        </div>
      </div>

      {/* Coluna Direita: Resultados */}
      <div className="w-full lg:w-7/12 flex flex-col">
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden flex flex-col h-[calc(100vh-8rem)]">
          
          <div className="px-5 py-4 border-b border-gray-100 bg-gray-50/50 flex justify-between items-center shrink-0">
            <h2 className="text-sm font-semibold text-gray-800">Resultado ({destVendor})</h2>
            <div className="flex items-center gap-4">
              {stats && (
                <span className="text-xs font-medium text-green-600 bg-green-50 px-2.5 py-1 rounded-full border border-green-100">
                  {stats.totalONUsFound} ONUs convertidas
                </span>
              )}
              {commands && (
                <div className="flex gap-2">
                  <button 
                    onClick={handleDownload}
                    className="text-xs font-medium text-blue-700 hover:text-blue-900 flex items-center gap-1.5 px-2.5 py-1 border border-blue-200 bg-blue-50 rounded hover:bg-blue-100 transition-colors shadow-sm"
                  >
                    <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
                    Baixar Arquivo (.txt)
                  </button>
                </div>
              )}
            </div>
          </div>

          {commands && (
            <div className="bg-gray-100/70 border-b border-gray-200 px-4 py-2 flex items-center justify-between gap-2 text-xs">
              <div className="flex gap-1 bg-gray-200/60 p-0.5 rounded-lg flex-wrap">
                <button
                  onClick={() => setViewMode("all")}
                  className={`px-3 py-1 rounded-md font-medium transition-all ${viewMode === "all" ? "bg-white text-gray-900 shadow-sm" : "text-gray-600 hover:text-gray-900"}`}
                >
                  📄 Script Completo
                </button>
                <button
                  onClick={() => setViewMode("clean")}
                  className={`px-3 py-1 rounded-md font-medium transition-all ${viewMode === "clean" ? "bg-white text-emerald-700 shadow-sm" : "text-gray-600 hover:text-gray-900"}`}
                >
                  🧹 Script Limpo (Sem Marcações)
                </button>
                <button
                  onClick={() => setViewMode("pon")}
                  className={`px-3 py-1 rounded-md font-medium transition-all ${viewMode === "pon" ? "bg-white text-blue-700 shadow-sm" : "text-gray-600 hover:text-gray-900"}`}
                >
                  ⚡ Cópia Rápida por PON
                </button>
                <button
                  onClick={() => setViewMode("onu")}
                  className={`px-3 py-1 rounded-md font-medium transition-all ${viewMode === "onu" ? "bg-white text-purple-700 shadow-sm" : "text-gray-600 hover:text-gray-900"}`}
                >
                  🔍 Cliente Individual
                </button>
              </div>
            </div>
          )}
          
          <div className="flex-1 bg-white relative p-0 overflow-hidden">
            {!commands ? (
              <div className="w-full h-full flex flex-col items-center justify-center text-gray-400 p-8 text-center bg-gray-50/30">
                <svg className="w-10 h-10 mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                <p className="text-sm font-medium">Os scripts de configuração aparecerão aqui.</p>
                <p className="text-xs mt-1 text-gray-400">Faça o upload do backup e execute a conversão.</p>
              </div>
            ) : viewMode === "all" ? (
              <div className="w-full h-full flex flex-col">
                <div className="p-2.5 bg-gray-50 border-b border-gray-100 flex justify-between items-center px-4">
                  <span className="text-xs text-gray-500 font-medium">Exibindo o script de migração completo com todas as portas PON:</span>
                  <button 
                    onClick={() => {navigator.clipboard.writeText(commands); alert("Script completo copiado!");}}
                    className="px-3 py-1 bg-gray-800 hover:bg-gray-900 text-white text-xs font-semibold rounded shadow-sm flex items-center gap-1.5 transition-all"
                  >
                    <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>
                    Copiar Script Completo
                  </button>
                </div>
                <textarea 
                  readOnly 
                  value={commands}
                  className="w-full flex-1 bg-transparent text-gray-800 font-mono text-[13px] p-5 focus:outline-none resize-none leading-relaxed"
                  spellCheck="false"
                ></textarea>
              </div>
            ) : viewMode === "clean" ? (
              <div className="w-full h-full flex flex-col">
                <div className="p-2.5 bg-emerald-50/50 border-b border-emerald-100 flex justify-between items-center px-4">
                  <span className="text-xs text-emerald-900 font-medium flex items-center gap-1.5">
                    <span className="w-2 h-2 rounded-full bg-emerald-500 inline-block"></span>
                    Exibindo o script limpo (sem as marcações por ONU/ONT para cópia contínua na OLT):
                  </span>
                  <button 
                    onClick={() => {navigator.clipboard.writeText(getCleanScript(commands)); alert("Script limpo copiado!");}}
                    className="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold rounded shadow flex items-center gap-1.5 transition-all"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>
                    Copiar Script Limpo
                  </button>
                </div>
                <textarea 
                  readOnly 
                  value={getCleanScript(commands)}
                  className="w-full flex-1 bg-transparent text-gray-800 font-mono text-[13px] p-5 focus:outline-none resize-none leading-relaxed"
                  spellCheck="false"
                ></textarea>
              </div>
            ) : viewMode === "pon" ? (
              <div className="w-full h-full flex flex-col">
                <div className="p-3 border-b border-gray-100 bg-blue-50/30 flex flex-wrap gap-2 items-center px-4">
                  <span className="text-xs font-bold text-blue-900 mr-2">Selecione a PON:</span>
                  {Object.keys(ponCommands).map((pon) => (
                    <button
                      key={pon}
                      onClick={() => setSelectedPon(pon)}
                      className={`px-3 py-1 text-xs rounded-lg font-semibold transition-all ${selectedPon === pon ? "bg-blue-600 text-white shadow" : "bg-white text-gray-700 border border-gray-300 hover:bg-gray-50"}`}
                    >
                      {pon}
                    </button>
                  ))}
                </div>
                <div className="p-2.5 bg-gray-50 border-b border-gray-100 flex justify-between items-center px-4">
                  <div className="flex items-center gap-3">
                    <span className="text-xs text-gray-600 font-medium">Bloco de comandos autônomo da PON:</span>
                    <label className="flex items-center gap-1.5 text-xs text-gray-700 font-semibold cursor-pointer bg-white px-2.5 py-1 rounded-lg border border-gray-200 shadow-2xs hover:bg-gray-50 transition-colors">
                      <input
                        type="checkbox"
                        checked={hideMarkers}
                        onChange={(e) => setHideMarkers(e.target.checked)}
                        className="rounded border-gray-300 text-blue-600 focus:ring-blue-500 w-3.5 h-3.5"
                      />
                      Sem marcações ONU/ONT
                    </label>
                  </div>
                  <button
                    onClick={() => { navigator.clipboard.writeText(hideMarkers ? getCleanScript(ponCommands[selectedPon] || "") : (ponCommands[selectedPon] || "")); alert(`Comandos da ${selectedPon} copiados!`); }}
                    className="px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white text-xs font-semibold rounded-lg shadow flex items-center gap-1.5 transition-all"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
                    Copiar Bloco da {selectedPon}
                  </button>
                </div>
                <textarea
                  readOnly
                  value={hideMarkers ? getCleanScript(ponCommands[selectedPon] || "") : (ponCommands[selectedPon] || "")}
                  className="w-full flex-1 bg-transparent text-gray-800 font-mono text-[13px] p-5 focus:outline-none resize-none leading-relaxed"
                  spellCheck="false"
                ></textarea>
              </div>
            ) : (
              <div className="w-full h-full flex flex-col">
                <div className="p-3 border-b border-gray-100 bg-purple-50/30 flex items-center gap-3 px-4">
                  <span className="text-xs font-bold text-purple-900 shrink-0">Buscar Cliente:</span>
                  <input
                    type="text"
                    placeholder="Digite nome, serial, slot ou ID (ex: maria, FHTT, 1/2/3)..."
                    value={searchOnu}
                    onChange={(e) => setSearchOnu(e.target.value)}
                    className="flex-1 px-3 py-1.5 text-xs border border-purple-200 rounded-lg focus:border-purple-500 focus:ring-1 focus:ring-purple-500 outline-none bg-white"
                  />
                </div>
                <div className="flex-1 overflow-y-auto p-4 space-y-3 bg-gray-50/50">
                  {Object.keys(onuCommands)
                    .filter((k) => k.toLowerCase().includes(searchOnu.toLowerCase()))
                    .map((onuKey) => (
                      <div key={onuKey} className="bg-white border border-gray-200 rounded-xl p-3.5 shadow-sm flex flex-col gap-2 hover:border-purple-200 transition-colors">
                        <div className="flex justify-between items-center border-b border-gray-100 pb-2">
                          <span className="text-xs font-bold text-gray-800 flex items-center gap-2">
                            <span className="w-2 h-2 rounded-full bg-purple-500"></span>
                            {onuKey}
                          </span>
                          <button
                            onClick={() => { navigator.clipboard.writeText(onuCommands[onuKey] || ""); alert(`Comandos copiados: ${onuKey}`); }}
                            className="px-2.5 py-1 bg-purple-600 hover:bg-purple-700 text-white text-xs font-semibold rounded-lg shadow-sm flex items-center gap-1.5 transition-all"
                          >
                            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>
                            Copiar SÓ este Cliente
                          </button>
                        </div>
                        <pre className="text-[11px] font-mono text-gray-700 bg-gray-50 p-2.5 rounded-lg border border-gray-150 overflow-x-auto whitespace-pre-wrap max-h-40 overflow-y-auto">{onuCommands[onuKey]}</pre>
                      </div>
                    ))}
                  {Object.keys(onuCommands).filter((k) => k.toLowerCase().includes(searchOnu.toLowerCase())).length === 0 && (
                    <div className="text-center py-8 text-gray-400 text-xs">Nenhum cliente encontrado para a busca "{searchOnu}".</div>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

    </div>
  );
}
