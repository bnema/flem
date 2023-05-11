import React, { useCallback } from 'react'
import { StyleSheet, View, Image, Text, ImageSourcePropType, Animated } from 'react-native'
import { CARD, ACTION_OFFSET } from '../../utils/Params'
import Selection from '../selection/Selection'

type MoviesProps = {
    name: string,
    source: ImageSourcePropType,
    isFirst: boolean,
    swipe: any, // any ?
    tiltSign: any, // any ?
}

export default function Card({ name, source, isFirst, swipe, tiltSign, ...rest }: MoviesProps) {

    // card move & rotation & effect :
    // --> multiply creates a new Animated value composed from two Animated values multiplied together.
    // --> interpolate allows input ranges to map to different output ranges.
    const rotate = Animated.multiply(swipe.x, tiltSign).interpolate({
        inputRange: [-ACTION_OFFSET, 0, ACTION_OFFSET],
        outputRange: ['-8deg', '0deg', '8deg'],
    })

    const yupOpacity = swipe.x.interpolate({
        inputRange: [20, ACTION_OFFSET],
        outputRange: [0, 1],
        extrapolate: 'clamp',
    })

    const nopeOpacity = swipe.x.interpolate({
        inputRange: [-ACTION_OFFSET, -20],
        outputRange: [1, 0],
        extrapolate: 'clamp',
    })

    // render of YUP or NOPE choices
    const renderSelection = useCallback(() => {
        return (
            <>
                <Animated.View style={[styles.selectionContainer, styles.yupContainer, { opacity: yupOpacity }]} >
                    <Selection type='YUP' />
                </Animated.View>
                <Animated.View style={[styles.selectionContainer, styles.nopeContainer, { opacity: nopeOpacity }]} >
                    <Selection type='NOPE' />
                </Animated.View>
            </>
        )
    }, [])

    const animatedCardStyle = {
        transform: [...swipe.getTranslateTransform(), { rotate }]
    }

    return (
        <Animated.View style={[styles.container, isFirst && animatedCardStyle]} {...rest}>
            <Image source={source} style={styles.image} />
            <Text style={styles.name} >{name}</Text>
            {isFirst && renderSelection()}
        </Animated.View>
    )
}

const styles = StyleSheet.create({
    container: {
        position: 'absolute',
        top: 45, // check for other device ?
    },
    image: {
        width: CARD.WIDTH,
        height: CARD.HEIGHT,
        borderRadius: CARD.BORDER_RADIUS,
    },
    name: {
        position: 'absolute',
        left: 22, // check for other device ?
        bottom: 22, // check for other device ?
        fontSize: 36,
        fontWeight: 'bold',
        color: '#fff',

    },
    selectionContainer: {
        position: 'absolute',
        top: 100, // check for other device ?
    },
    yupContainer: {
        left: 45, // check for other device ?
        transform: [{ rotate: '-30deg' }]
    },
    nopeContainer: {
        right: 45, // check for other device ?
        transform: [{ rotate: '30deg' }]
    },
})