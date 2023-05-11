import React, { useState, useRef, useCallback, useEffect } from 'react'
import axios from 'axios'
import { StyleSheet, View, Animated, PanResponder } from 'react-native'

import Card from '../card/Card'
import { CARD, ACTION_OFFSET } from '../../utils/Params'
import { nanars as nanarsArray } from './data'

export default function Swiper() {
    const [nanars, setNanars] = useState(nanarsArray)
    // const [moviesData, setMoviesData] = useState([])

    const swipe = useRef(new Animated.ValueXY()).current
    const tiltSign = useRef(new Animated.Value(1)).current

    useEffect(() => {
        fetch(`http://192.168.1.151:3333/v1/tmdb/random10`)
        .then(response => response.json())
        .then(data => {
            const caca = data
            console.log('movies data', data)
        }).catch((err) => {
            console.log('fetch error', err);     
        })
    }, [])

    // gesture event parameters (card moves on touch):
    const panResponder = PanResponder.create({
        onMoveShouldSetPanResponder: () => true,
        onPanResponderMove: (_, { dx, dy, y0 }) => {
            swipe.setValue({ x: dx, y: dy })
            tiltSign.setValue(y0 > CARD.HEIGHT / 2 ? 1 : 0.4) // adjusts the amplitude of tilt depending on whether the card is moved from the bottom or the top.
        },
        onPanResponderRelease: (_, { dx, dy }) => {
            const direction = Math.sign(dx)
            const isActive = Math.abs(dx) > ACTION_OFFSET
            if (isActive) {
                Animated.timing(swipe, { // animates a value along a timed easing curve.
                    duration: 300,
                    toValue: {
                        x: direction * 500,
                        y: dy,
                    },
                    useNativeDriver: true // send everything about the animation to native before starting the animation, allowing native code to perform the animation on the UI thread without having to go through the bridge on every frame.
                }).start(removeTopCard) // animations are started by calling start().

            } else {
                Animated.spring(swipe, { // animates a value, tracks velocity state to create fluid motions as the toValue updates.
                    toValue: {
                        x: 0,
                        y: 0,
                    },
                    useNativeDriver: true,
                    friction: 5, // bounciness
                }).start()
            }
        }
    })

    // remove the top card after swipe
    const removeTopCard = useCallback(() => {
        setNanars((prevState) => {
            console.log("🚀 ~ file: Swiper.tsx:65 ~ removeTopCard ~ prevState:", prevState[0]) // 
            return prevState.slice(1)
        })
        swipe.setValue({ x: 0, y: 0 })
    }, [swipe])

    return (
        <View style={styles.container} >
            {nanars.map(({ name, source }, index) => {
                const isFirst = index === 0;
                const dragHandlers = isFirst ? panResponder.panHandlers : {};

                return (<Card key={name} name={name} source={source} isFirst={isFirst} swipe={swipe} tiltSign={tiltSign} {...dragHandlers} />)
            }).reverse()}
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        alignItems: 'center',
    }
})