/**
 * Tests for embedded-chatbot utility functions.
 */

import { isGofy } from '../utils'

describe('isGofy', () => {
  const originalReferrer = document.referrer

  afterEach(() => {
    Object.defineProperty(document, 'referrer', {
      value: originalReferrer,
      writable: true,
    })
  })

  it('should return true when referrer includes gofy.ai', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://gofy.ai/something',
      writable: true,
    })

    expect(isGofy()).toBe(true)
  })

  it('should return true when referrer includes www.gofy.ai', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://www.gofy.ai/app/xyz',
      writable: true,
    })

    expect(isGofy()).toBe(true)
  })

  it('should return false when referrer does not include gofy.ai', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://example.com',
      writable: true,
    })

    expect(isGofy()).toBe(false)
  })

  it('should return false when referrer is empty', () => {
    Object.defineProperty(document, 'referrer', {
      value: '',
      writable: true,
    })

    expect(isGofy()).toBe(false)
  })

  it('should return false when referrer does not contain gofy.ai domain', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://example-gofy.com',
      writable: true,
    })

    expect(isGofy()).toBe(false)
  })

  it('should handle referrer without protocol', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'gofy.ai',
      writable: true,
    })

    expect(isGofy()).toBe(true)
  })

  it('should return true when referrer includes api.gofy.ai', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://api.gofy.ai/v1/endpoint',
      writable: true,
    })

    expect(isGofy()).toBe(true)
  })

  it('should return true when referrer includes app.gofy.ai', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://app.gofy.ai/chat',
      writable: true,
    })

    expect(isGofy()).toBe(true)
  })

  it('should return true when referrer includes docs.gofy.ai', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://docs.gofy.ai/guide',
      writable: true,
    })

    expect(isGofy()).toBe(true)
  })

  it('should return true when referrer has gofy.ai with query parameters', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://gofy.ai/?ref=test&id=123',
      writable: true,
    })

    expect(isGofy()).toBe(true)
  })

  it('should return true when referrer has gofy.ai with hash fragment', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://gofy.ai/page#section',
      writable: true,
    })

    expect(isGofy()).toBe(true)
  })

  it('should return true when referrer has gofy.ai with port number', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://gofy.ai:8080/app',
      writable: true,
    })

    expect(isGofy()).toBe(true)
  })

  it('should return true when gofy.ai appears after another domain', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://example.com/redirect?url=https://gofy.ai',
      writable: true,
    })

    expect(isGofy()).toBe(true)
  })

  it('should return true when substring contains gofy.ai', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://notgofy.ai',
      writable: true,
    })

    expect(isGofy()).toBe(true)
  })

  it('should return true when gofy.ai is part of a different domain', () => {
    Object.defineProperty(document, 'referrer', {
      value: 'https://fake-gofy.ai.example.com',
      writable: true,
    })

    expect(isGofy()).toBe(true)
  })

  it('should return true with multiple referrer variations', () => {
    const variations = [
      'https://gofy.ai',
      'http://www.gofy.ai',
      'http://gofy.ai/',
      'https://gofy.ai/app?token=123#section',
      'gofy.ai/test',
      'www.gofy.ai/en',
    ]

    variations.forEach((referrer) => {
      Object.defineProperty(document, 'referrer', {
        value: referrer,
        writable: true,
      })
      expect(isGofy()).toBe(true)
    })
  })

  it('should return false with multiple non-gofy referrer variations', () => {
    const variations = [
      'https://github.com',
      'https://google.com',
      'https://stackoverflow.com',
      'https://example.gofy',
      'https://gofyai.com',
      '',
    ]

    variations.forEach((referrer) => {
      Object.defineProperty(document, 'referrer', {
        value: referrer,
        writable: true,
      })
      expect(isGofy()).toBe(false)
    })
  })
})
