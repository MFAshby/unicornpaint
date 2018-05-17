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

function getPixel(x, y, pixels) {
    let row = pixels[y]
    if (!row) {
        return
    }
    return row[x]
}

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
  

export {
    xy,
    rgb,
    colorEqual,
    coordsEqual,
    findContiguousPixels,
    getPixel
}

