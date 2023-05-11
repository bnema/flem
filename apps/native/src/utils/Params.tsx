import { Dimensions } from 'react-native'

const { height, width } = Dimensions.get('screen')

export const CARD = {
    WIDTH: width * .9,
    HEIGHT: height * .78,
    BORDER_RADIUS: 20,
}

export const COLORS = {
    YUP: '#00ffb7',
    NOPE: '#ff195e',
}

export const ACTION_OFFSET = 100