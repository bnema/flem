import React, { useCallback, useRef, useState, useEffect } from 'react'
import { StyleSheet, Text, Animated, Pressable } from 'react-native'
import { COLORS } from '../../utils/Params'
import colorsDB from '../../utils/colorsDB'

type TypeProps = {
    id: number,
    genre: string
}

export default function GenreSelection({ id, genre }: TypeProps) {

    const [colordButton, setColorButton] = useState(false)
    const [newColor, setNewColor] = useState('')
    const [genreSelection, setGenreSelection] = useState<{ id: number, genre: string }[]>([])
    console.log("🚀 ~ file: GenreSelection.tsx:16 ~ GenreSelection ~ genreSelection:", genreSelection)

    const scale = useRef(new Animated.Value(1)).current

    const AnimatedScale = useCallback((newValue: number) => {
        Animated.spring(scale, {
            toValue: newValue,
            friction: 4,
            useNativeDriver: true,
        }).start()
    }, [scale])

    // set choices : genres, button random color
    const handleSelection = () => {
        setColorButton(!colordButton)
        !colordButton ? setGenreSelection([...genreSelection, { id: id, genre: genre }]) : setGenreSelection([...genreSelection.filter(genre => genre.id != id)])

        // generate a random color attribution
        const randomColor = colorsDB[Math.floor(Math.random() * colorsDB.length)]
        setNewColor(randomColor)
    }

    return (
        <Pressable onPressIn={() => { handleSelection(); AnimatedScale(.8) }} onPressOut={() => AnimatedScale(1)} id={`${id}`}>
            {!colordButton
                ? <Animated.View style={[styles.container, { transform: [{ scale }] }, { backgroundColor: COLORS.BACKGROUND }]} >
                    <Text style={[styles.genre, { color: COLORS.VIEW }]}>{genre}</Text>
                </Animated.View>
                : <Animated.View style={[styles.container, { transform: [{ scale }] }, { backgroundColor: newColor }, { borderColor: COLORS.UI }]} >
                    <Text style={[styles.genre, { color: COLORS.BACKGROUND }]}>{genre}</Text>
                </Animated.View>
            }
        </Pressable>
    )
}

const styles = StyleSheet.create({
    container: {
        borderWidth: 1,
        borderRadius: 5,
        borderColor: COLORS.VIEW,
        padding: 16,
    },
    genre: {
        fontWeight: '500',
    }
})