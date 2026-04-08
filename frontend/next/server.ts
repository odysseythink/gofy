export class NextResponse {
  static json(data: any, init?: ResponseInit) {
    return new Response(JSON.stringify(data), {
      ...init,
      headers: { 'Content-Type': 'application/json', ...init?.headers },
    })
  }

  static redirect(url: string | URL, status = 307) {
    return Response.redirect(url, status)
  }
}

export type NextRequest = Request
