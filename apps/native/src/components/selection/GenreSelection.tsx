import React, { useCallback, useRef, useState } from 'react'
import { StyleSheet, Text, Animated, Pressable } from 'react-native'
import { COLORS } from '../../utils/Params'
import colorsDB from '../../utils/colorsDB'

type SelectionProps = {
    id: string,
    genre: string
}

export default function GenreSelection({ id, genre }: SelectionProps) {

    const [buttonColor, setButtonColor] = useState(false)
    const [newColor, setNewColor] = useState('')

    const scale = useRef(new Animated.Value(1)).current

    const AnimatedScale = useCallback((newValue) => {
        Animated.spring(scale, {
            toValue: newValue,
            friction: 4,
            useNativeDriver: true,
        }).start()
    }, [scale])

    // generate a random color attribution onPress
    const colorGenerator = () => {
        setButtonColor(!buttonColor)
        const randomColor = colorsDB[Math.floor(Math.random() * colorsDB.length)]
        setNewColor(randomColor)
    }

    return (
        <Pressable onPressIn={() => { colorGenerator(); AnimatedScale(.8) }} onPressOut={() => AnimatedScale(1)}>
            {!buttonColor
                ? <Animated.View style={[styles.container, { transform: [{ scale }] }, { backgroundColor: COLORS.BACKGROUND }]} id={id}>
                    <Text style={[styles.genre, { color: COLORS.VIEW }]}>{genre}</Text>
                </Animated.View>
                : <Animated.View style={[styles.container, { transform: [{ scale }] }, { backgroundColor: newColor }]} id={id}>
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