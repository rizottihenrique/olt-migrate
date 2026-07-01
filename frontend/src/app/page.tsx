"use client";

import { useState } from "react";

export default function Home() {
  const [file, setFile] = useState<File | null>(null);
  const [sourceVendor, setSourceVendor] = useState("Fiberhome AN5000");
  const [destVendor, setDestVendor] = useState("Nokia");
  const [mappings, setMappings] = useState([{ sourceSlot: 1, sourcePon: 1, destSlot: 1, destPon: 1 }]);
  const [commands, setCommands] = useState("");
  const [loading, setLoading] = useState(false);
  const [stats, setStats] = useState<{ totalONUsFound: number } | null>(null);

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
    if (!file) return alert("Selecione um arquivo de configuração da Fiberhome.");
    
    setLoading(true);
    const formData = new FormData();
    formData.append("configFile", file);
    formData.append("sourceVendor", sourceVendor);
    formData.append("destVendor", destVendor);
    formData.append("mappings", JSON.stringify(mappings));

    try {
      const res = await fetch("http://localhost:8080/api/migrate", {
        method: "POST",
        body: formData,
      });

      if (!res.ok) {
        throw new Error("Erro na migração: " + await res.text());
      }

      const data = await res.json();
      setCommands(data.commands);
      setStats({ totalONUsFound: data.totalONUsFound });
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

        {/* Passo 2: Arquivo */}
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
          <div className="px-5 py-4 border-b border-gray-100 bg-gray-50/50">
            <h2 className="text-sm font-semibold text-gray-800">2. Arquivo de Backup</h2>
          </div>
          <div className="p-5">
            <label className="block w-full border-2 border-dashed border-gray-200 hover:border-blue-400 hover:bg-blue-50/50 transition-colors rounded-lg p-6 text-center cursor-pointer">
              <input type="file" onChange={handleFileChange} className="hidden" accept=".txt,.dat,.cfg" />
              {file ? (
                <div>
                  <svg className="w-8 h-8 mx-auto mb-2 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>
                  <p className="text-sm font-medium text-gray-800">{file.name}</p>
                  <p className="text-xs text-gray-500 mt-1">{(file.size / 1024).toFixed(1)} KB</p>
                </div>
              ) : (
                <div>
                  <svg className="w-8 h-8 mx-auto mb-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
                  <p className="text-sm font-medium text-gray-600">Selecione o backup Fiberhome</p>
                </div>
              )}
            </label>
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
              disabled={loading || !file}
              className={`w-full py-2.5 rounded-lg text-sm font-medium text-white shadow-sm transition-all mt-4 flex items-center justify-center
                ${loading || !file ? 'bg-gray-300 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700 hover:shadow'}`}
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
            <h2 className="text-sm font-semibold text-gray-800">Resultado (Comandos Nokia)</h2>
            <div className="flex items-center gap-4">
              {stats && (
                <span className="text-xs font-medium text-green-600 bg-green-50 px-2.5 py-1 rounded-full border border-green-100">
                  {stats.totalONUsFound} ONUs convertidas
                </span>
              )}
              {commands && (
                <div className="flex gap-2">
                  <button 
                    onClick={() => {navigator.clipboard.writeText(commands); alert("Copiado!");}}
                    className="text-xs font-medium text-gray-600 hover:text-gray-900 flex items-center gap-1.5 px-2.5 py-1 border border-gray-200 rounded hover:bg-gray-50 transition-colors"
                  >
                    <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>
                    Copiar
                  </button>
                  <button 
                    onClick={handleDownload}
                    className="text-xs font-medium text-blue-700 hover:text-blue-900 flex items-center gap-1.5 px-2.5 py-1 border border-blue-200 bg-blue-50 rounded hover:bg-blue-100 transition-colors shadow-sm"
                  >
                    <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
                    Baixar .txt
                  </button>
                </div>
              )}
            </div>
          </div>
          
          <div className="flex-1 bg-white relative p-1 overflow-hidden">
            {commands ? (
              <textarea 
                readOnly 
                value={commands}
                className="w-full h-full bg-transparent text-gray-800 font-mono text-[13px] p-5 focus:outline-none resize-none leading-relaxed"
                spellCheck="false"
              ></textarea>
            ) : (
              <div className="w-full h-full flex flex-col items-center justify-center text-gray-400 p-8 text-center bg-gray-50/30">
                <svg className="w-10 h-10 mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                <p className="text-sm font-medium">Os scripts de configuração aparecerão aqui.</p>
                <p className="text-xs mt-1 text-gray-400">Faça o upload do backup e execute a conversão.</p>
              </div>
            )}
          </div>
        </div>
      </div>

    </div>
  );
}
