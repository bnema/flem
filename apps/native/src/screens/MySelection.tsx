import { StyleSheet, View } from "react-native"
import { COLORS } from "../utils/Params"

export default function MySelection() {
    return (
        <View style={styles.container}>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        // flexDirection: "column",
        // justifyContent: 'center',
        // alignItems: 'center',
        backgroundColor: COLORS.BACKGROUND,
    }
})