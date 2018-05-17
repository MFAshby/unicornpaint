import { colorEqual, coordsEqual, xy, rgb, findContiguousPixels, rotatePixelsClock, rotatePixelsCounterClock } from './Utils'

test('test colorEqual function', () => {
    expect(colorEqual([0, 0, 0], [0, 0, 0]))
        .toBe(true)
    expect(colorEqual([1, 0, 0], [0, 0, 0]))
        .toBe(false)
    expect(colorEqual([0, 1, 0], [0, 0, 0]))
        .toBe(false)
    expect(colorEqual([0, 0, 1], [0, 0, 0]))
        .toBe(false)

    expect(colorEqual([0, 0, 1], {r: 0, g: 0, b: 1}))
        .toBe(true)
    expect(colorEqual({r: 0, g: 0, b: 1}, [0, 0, 1]))
        .toBe(true)

    expect(colorEqual({r: 0, g: 0, b: 0}, [0, 0, 1]))
        .toBe(false)
    expect(colorEqual([0, 0, 1], {r: 0, g: 0, b: 0}))
        .toBe(false)

    expect(colorEqual(undefined, {r: 0, g: 0, b: 0}))
        .toBe(false)
    expect(colorEqual({r: 0, g: 0, b: 0}, undefined))
        .toBe(false)
})

test('coordEqual function', () => {
    expect(coordsEqual([0, 0], [0, 0]))
        .toBe(true)
    expect(coordsEqual([0, 1], [0, 0]))
        .toBe(false)
    expect(coordsEqual([1, 0], [0, 0]))
        .toBe(false)

    expect(coordsEqual([1, 0], [0, 0]))
        .toBe(false)

    expect(coordsEqual(undefined, [0, 0]))
        .toBe(false)
    expect(coordsEqual([0, 0], undefined))
        .toBe(false)
})

test('contiguousPixels', () => {
    let allBlack = [
        [[0, 0, 0], [0, 0, 0], [0, 0, 0], [0, 0, 0]],
        [[0, 0, 0], [0, 0, 0], [0, 0, 0], [0, 0, 0]],
        [[0, 0, 0], [0, 0, 0], [0, 0, 0], [0, 0, 0]],
        [[0, 0, 0], [0, 0, 0], [0, 0, 0], [0, 0, 0]]
    ]

    let contiguousPx = findContiguousPixels(0, 0, allBlack)
    expect(contiguousPx.length)
        .toBe(16)

    expect(findContiguousPixels(0, 0, [
        [[0, 0, 0], [0, 0, 1], [0, 0, 0], [0, 0, 0]],
        [[0, 0, 1], [0, 0, 1], [0, 0, 0], [0, 0, 0]],
        [[0, 0, 0], [0, 0, 0], [0, 0, 0], [0, 0, 0]],
        [[0, 0, 0], [0, 0, 0], [0, 0, 0], [0, 0, 0]]
    ]).length)
        .toBe(1)

    expect(findContiguousPixels(1, 1, [
        [[0, 0, 0], [0, 0, 0], [0, 0, 0], [0, 0, 0]],
        [[0, 0, 0], [0, 0, 1], [0, 0, 1], [0, 0, 0]],
        [[0, 0, 0], [0, 0, 1], [0, 0, 1], [0, 0, 0]],
        [[0, 0, 0], [0, 0, 0], [0, 0, 0], [0, 0, 0]]
    ]).length)
        .toBe(4)

    expect(findContiguousPixels(1, 1, [
        [[0, 0, 1], [0, 0, 1], [0, 0, 1], [0, 0, 1]],
        [[0, 0, 1], [0, 0, 0], [0, 0, 0], [0, 0, 1]],
        [[0, 0, 1], [0, 0, 0], [0, 0, 0], [0, 0, 1]],
        [[0, 0, 1], [0, 0, 1], [0, 0, 1], [0, 0, 1]]
    ]).length)
        .toBe(4)
})

test('rotate clockwise function', () => {
    let rot1 = rotatePixelsClock([
        [[0, 0, 1], [0, 0, 0], [0, 0, 0], [0, 0, 2]],
        [[0, 0, 0], [0, 0, 0], [0, 0, 0], [0, 0, 0]],
        [[0, 0, 0], [0, 0, 0], [0, 0, 0], [0, 0, 0]],
        [[0, 0, 4], [0, 0, 0], [0, 0, 0], [0, 0, 3]]
    ])
    expect(rot1[0][3][2])
    .toBe(1)
    expect(rot1[3][3][2])
    .toBe(2)
    expect(rot1[3][0][2])
    .toBe(3)
    expect(rot1[0][0][2])
    .toBe(4)
})

test('rotate counter-clockwise function', () => {
    let rot1 = rotatePixelsCounterClock([
        [[0, 0, 1], [0, 0, 0], [0, 0, 0], [0, 0, 2]],
        [[0, 0, 0], [0, 0, 0], [0, 0, 0], [0, 0, 0]],
        [[0, 0, 0], [0, 0, 0], [0, 0, 0], [0, 0, 0]],
        [[0, 0, 4], [0, 0, 0], [0, 0, 0], [0, 0, 3]]
    ])
    expect(rot1[0][3][2])
    .toBe(3)
    expect(rot1[3][3][2])
    .toBe(4)
    expect(rot1[3][0][2])
    .toBe(1)
    expect(rot1[0][0][2])
    .toBe(2)
})