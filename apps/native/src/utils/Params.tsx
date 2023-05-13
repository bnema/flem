import { Dimensions } from 'react-native'

const { height, width } = Dimensions.get('screen')

export const CARD = {
    WIDTH: width * .9,
    HEIGHT: height * .78,
    BORDER_RADIUS: 20,
}

export const COLORS = {
    BACKGROUND: '#25224f',
    UI: '#15132e',
    ICON: '#a8a0ff',
    BUTTON: '#9288fc',
    YUP: '#00ffb7',
    NOPE: '#ff195e',
    VIEW: '#d6d1ff',
}

export const ACTION_OFFSET = 100