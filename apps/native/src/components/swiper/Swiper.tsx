import React, { useState, useRef, useCallback, useEffect } from 'react'
import { StyleSheet, View, Animated, PanResponder } from 'react-native'
import Card from '../card/Card'
import { CARD, ACTION_OFFSET } from '../../utils/Params'
import { Movie } from "@flem/types"

export default function Swiper() {
    const [moviesData, setMoviesData] = useState([])

    const swipe = useRef(new Animated.ValueXY()).current
    const tiltSign = useRef(new Animated.Value(1)).current

    useEffect(() => {
        fetch(`http://api.flem.bnema.dev/v1/tmdb/movies?genre=28&minDate=2000-01-01&maxDate=2010-12-31&quantity=5`)
            .then(response => response.json())
            .then(data => {

                const movies: Movie[] = data
                const moviesData = movies.map(movie => {
                    const poster = `https://image.tmdb.org/t/p/w500/${movie.poster_path}` //path for poster
                    const genre = movie.genres.map(genre => genre.name) // extract 'genres' and put them in array
                    const genreID = movie.genres.map(genre => genre.id) // extract 'genres ids' and put them in array
                    return { id: movie.id, title: movie.title, genre, genreID, overview: movie.overview, date: movie.release_date, poster }
                })

                setMoviesData(moviesData)

            }).catch((error) => {
                console.log('fetch API error', error);
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
            const directionX = Math.sign(dx)
            const directionY = Math.sign(dy)
            const isActiveX = Math.abs(dx) > ACTION_OFFSET
            const isActiveY = dy < -ACTION_OFFSET

            if (isActiveX) {
                Animated.timing(swipe, { // animates a value along a timed easing curve.
                    duration: 300,
                    toValue: {
                        x: directionX * 500,
                        y: dy,
                    },
                    useNativeDriver: true // send everything about the animation to native before starting the animation, allowing native code to perform the animation on the UI thread without having to go through the bridge on every frame.
                }).start(removeTopCard) // animations are started by calling start().

            } else if (isActiveY) {
                Animated.timing(swipe, { // animates a value along a timed easing curve.
                    duration: 300,
                    toValue: {
                        x: 0,
                        y: -Math.abs(directionY * 260),
                    },
                    useNativeDriver: true // send everything about the animation to native before starting the animation, allowing native code to perform the animation on the UI thread without having to go through the bridge on every frame.
                }).start(viewTopCard)

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

    // remove the top card after swipe (& send infos to server???)
    const removeTopCard = useCallback(() => {
        setMoviesData((prevState) => {
            return prevState.slice(1)
        })

        const swipeXInfo = parseInt(JSON.stringify(swipe.x)) // convert type of Animated to number
        if (swipeXInfo > 0) {
            console.log('RIGHT is YUP')
        } else if (swipeXInfo < 0) {
            console.log('LEFT is NOPE')
        }

        swipe.setValue({ x: 0, y: 0 })
    }, [swipe])

    // show info of the overview of card 
    const viewTopCard = useCallback(() => {
        // setMoviesData((prevState) => {
        //     return prevState
        // })

        const swipeYInfo = parseInt(JSON.stringify(swipe.y)) // convert type of Animated to number

        if (swipeYInfo < 0) console.log('UP is for VIEW');

    }, [swipe])

    return (
        <View style={styles.container} >
            {moviesData.map(({ title, genre, overview, date, poster }, index) => {
                const isFirst = index === 0;
                const dragHandlers = isFirst ? panResponder.panHandlers : {};

                return (<Card key={title} title={title} genre={genre} overview={overview} date={date} poster={poster} isFirst={isFirst} swipe={swipe} tiltSign={tiltSign} {...dragHandlers} />)
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