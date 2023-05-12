import React, { useCallback, useState } from 'react'
import { StyleSheet, View, Image, Animated } from 'react-native'
import { CARD, ACTION_OFFSET } from '../../utils/Params'
import Selection from '../selection/Selection'
import CardOverview from './CardOverview'
import EyeButton from '../Button/EyeButton'

type TypeProps = {
    title: string,
    genre: string[],
    overview: string,
    date: string,
    poster: string,
    isFirst: boolean,
    swipe: any, //  ?
    tiltSign: any, //  ?
}

export default function Card({ title, genre, overview, date, poster, isFirst, swipe, tiltSign, ...rest }: TypeProps) {

    const [watched, setwatched] = useState(false)

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

    const viewOpacity = swipe.y.interpolate({
        inputRange: [-180, -20],
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
                <Animated.View style={[styles.selectionContainer, styles.viewContainer, { opacity: viewOpacity }]} >
                    <Selection type='VIEW' />
                </Animated.View>
            </>
        )
    }, [])

    const animatedCardStyle = {
        transform: [...swipe.getTranslateTransform(), { rotate }]
    }

    // display card's overview
    const overviewOpacity = swipe.y.interpolate({
        inputRange: [-260, -200],
        outputRange: [1, 0],
        extrapolate: 'clamp',
    })

    const renderOverview = useCallback(() => {
        return (
            <>
                <Animated.View style={[styles.cardOverviewContainer, { opacity: overviewOpacity }]} >
                    <CardOverview title={title} genre={genre} overview={overview} date={date} />
                </Animated.View>
            </>
        )
    }, [])


    return (
        <Animated.View style={[styles.container, isFirst && animatedCardStyle]} {...rest}>
            <Image source={{ uri: poster }} alt={title} style={styles.image} />

            {watched === false
                ? <EyeButton onPressIn={()=> setwatched(!watched)} name={'eye-slash'} size={50} color={'#ff195e'} />
                : <EyeButton onPressIn={()=>setwatched(!watched)} name={'eye'} size={50} color={'#00ffb7'} />
            }

            {isFirst && (
                <>
                    {renderSelection()}
                    {renderOverview()}
                </>
            )}
        </Animated.View>

    )
}

const styles = StyleSheet.create({
    container: {
        position: 'absolute',
        alignItems: 'center',
        top: 45, // check for other device ?
    },
    image: {
        width: CARD.WIDTH,
        height: CARD.HEIGHT,
        borderRadius: CARD.BORDER_RADIUS,
        resizeMode: 'contain',
    },
    title: {
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
    viewContainer: {
        top: 300, // check for other device ?
    },
    cardOverviewContainer: {
        position: 'relative',
        zIndex: -1,
        // top: -20,
        bottom: 70,
    },
})