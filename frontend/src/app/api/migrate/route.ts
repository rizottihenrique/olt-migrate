import { NextRequest, NextResponse } from "next/server";

export async function POST(req: NextRequest) {
  try {
    const contentType = req.headers.get("content-type") || "";
    const bodyBuffer = await req.arrayBuffer();
    
    // Conecta diretamente no Go usando localhost (IPv6 [::1] no Windows)
    const res = await fetch("http://localhost:8080/api/migrate", {
      method: "POST",
      headers: { "Content-Type": contentType },
      body: bodyBuffer,
    });

    if (!res.ok) {
      const errorText = await res.text();
      return NextResponse.json({ error: errorText }, { status: res.status });
    }

    const data = await res.json();
    return NextResponse.json(data);
  } catch (err: any) {
    console.error("Erro no proxy API Next.js:", err, err?.cause);
    return NextResponse.json(
      { error: "Erro ao conectar com o backend Go: " + (err.message || err) },
      { status: 500 }
    );
  }
}
