import base64js from 'base64-js'
import omggif from 'omggif'

/**
 * Converts RGB array to named object
 * @param {[]} item, a length 3 array containing 8 bit RGB value
 */
function rgb(item) {
    return {
        r: item[0],
        g: item[1],
        b: item[2]
    }
}

/**
 * Converts XY array to named object
 * @param {[]} item, a 2 length array containing x and y values
 */
function xy(item) {
    return {
        x: item[0],
        y: item[1]
    }
}

/**
 * Compares colors. 
 * Colors are either 3-element arrays [r, g, b]
 * Or objects with .r .g .b members
 */
function colorEqual(item1, item2) {
    if (!item1 || !item2) {
        return false
    }
    if (Array.isArray(item1)) {
        item1 = rgb(item1)
    }
    if (Array.isArray(item2)) {
        item2 = rgb(item2)
    }
    return item1.r === item2.r
        && item1.g === item2.g
        && item1.b === item2.b
}

/**
 * Compares 2 coordinate values
 * @param { [x, y] or {x: y:} } item1 
 * @param { [x, y] or {x: y:} } item2 
 */
function coordsEqual(item1, item2) {
    if (!item1 || !item2) {
        return false
    }

    if (Array.isArray(item1)) {
        item1 = xy(item1)
    }
    if (Array.isArray(item2)) {
        item2 = xy(item2)
    }
    return item1.x === item2.x 
        && item1.y === item2.y
}

/**
 * If I ever want to change the pixel format...
 * I should really route all access to pixels through here.
 * @param {number} x 
 * @param {number} y 
 * @param {[][][]} pixels 
 */
function getPixel(x, y, pixels) {
    let column = pixels[x]
    if (!column) {
        return
    }
    return column[y]
}

/**
 * Finds contiguous regions of pixels of the same color.
 * Returns an array of the x & y coordinates of every pixel in the region.
 * Doesn't jump corners: only pixels that share a side are considered to 
 * join up.
 * 
 * @param {number} x 
 * @param {number} y 
 * @param {[][]]} pixels 
 * @param {[]]} targetColor 
 * @param {[][]]} contiguousPixels 
 */
function findContiguousPixels(x, y, pixels, targetColor = getPixel(x, y, pixels), contiguousPixels=[[x, y]]) {
    let adjescent = [
      [x-1, y],
      [x+1, y],
      [x, y-1],
      [x, y+1]
    ]
  
    adjescent.forEach((coord) => {
      let px = xy(coord)
      let pxCol = getPixel(px.x, px.y, pixels)
      if (!pxCol) {
          return
      }

      // add adjescents uniquely if they are the target color
      let ix = contiguousPixels.findIndex((existingCoord) => coordsEqual(coord, existingCoord))
      if (ix !== -1) {
          return 
      }

      if (!colorEqual(pxCol, targetColor)) {
          return
      }
      contiguousPixels.push(coord)
      let morePixels = findContiguousPixels(px.x, px.y, pixels, targetColor, contiguousPixels)
      contiguousPixels.concat(morePixels)
    })
  
    return contiguousPixels
}

function rotatePixelsCounterClock(pixels) {
    let rotateClock = (x, y, width, height) => {
        return {
            newx: -y + (width - 1), 
            newy : x
        }
    }
    return transformPixels(pixels, rotateClock)
}

function rotatePixelsClock(pixels) {
    let rotateClock = (x, y, width, height) => {
        return {
            newx: y, 
            newy : - x + (height - 1)
        }
    }
    return transformPixels(pixels, rotateClock)
}

/**
 * Apply an arbitrary transform to some pixels.
 * Does not modify the original pixels, just returns the new ones
 * 
 * @param {[][][]]} pixels 
 * @param {(x, y, width, height) => { newx: newy: }} transform 
 * @returns {[][][]} newPixels
 */
function transformPixels(pixels, transform) {
    let width = pixels.length
    let height = pixels[0].length
    let newPixels = []
    for (let x = 0; x < width; x++) {
        let column = []
        for (var y = 0; y < height; y++) {
            column.push([0, 0, 0])
        }
        newPixels.push(column)
    }
    
    for (let x = 0; x < width; x++) {
      for (let y = 0; y < height; y++) {
        let px = getPixel(x, y, pixels)
        let {newx, newy} = transform(x, y, width, height)
        newPixels[newx][newy] = px
      }
    }
    return newPixels
}

/**
 * Reads all frames from a GIF image, returns them as 3d arrays (x, y, color component)
 * @param {string} base64GifData The GIF image encoded as base64
 */
function readGifFrames(base64GifData) {
    let imageBytes = base64js.toByteArray(base64GifData)
    let gifReader = new omggif.GifReader(imageBytes)
    
    let { width, height } = gifReader
    let numFrames = gifReader.numFrames()
    let frames = []

    for (var i=0; i<numFrames; i++) {
      let rawPixels = new Array(width * height * 4)
      gifReader.decodeAndBlitFrameRGBA(i, rawPixels)

      // Create the x, y array upfront
      let pixels = new Array(width)
      for (let y=0; y<height; y++) {
        pixels[y] = new Array(height)
      }
      frames.push(pixels)

      // Copy pixels to out array. The data provided is provided in rows
      var ix = 0
      for (let y=0; y<height; y++) {
        for (let x=0; x<width; x++) {
          let r = rawPixels[ix++]
          let g = rawPixels[ix++]
          let b = rawPixels[ix++]
          ix++ // Ignore the alpha component
          pixels[x][y] = [r, g, b]
        }
      }
    }

    return frames
}

export {
    xy,
    rgb,
    colorEqual,
    coordsEqual,
    findContiguousPixels,
    getPixel,
    rotatePixelsClock,
    rotatePixelsCounterClock,
    readGifFrames
}

